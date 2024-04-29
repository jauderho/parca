// Copyright 2022-2024 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"io"
	"os"
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/memory"
	"github.com/go-kit/log"
	pprofprofile "github.com/google/pprof/profile"
	columnstore "github.com/polarsignals/frostdb"
	"github.com/polarsignals/frostdb/query"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pprofpb "github.com/parca-dev/parca/gen/proto/go/google/pprof"
	profilestorepb "github.com/parca-dev/parca/gen/proto/go/parca/profilestore/v1alpha1"
	pb "github.com/parca-dev/parca/gen/proto/go/parca/query/v1alpha1"
	sharepb "github.com/parca-dev/parca/gen/proto/go/parca/share/v1alpha1"
	"github.com/parca-dev/parca/pkg/ingester"
	"github.com/parca-dev/parca/pkg/kv"
	"github.com/parca-dev/parca/pkg/normalizer"
	"github.com/parca-dev/parca/pkg/parcacol"
	"github.com/parca-dev/parca/pkg/profile"
	"github.com/parca-dev/parca/pkg/profilestore"
)

func getShareServerConn(t Testing) sharepb.ShareServiceClient {
	conn, err := grpc.Dial("api.pprof.me:443", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	require.NoError(t, err)
	return sharepb.NewShareServiceClient(conn)
}

func TestColumnQueryAPIQueryRangeEmpty(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	_, err = colDB.Table("stacktraces", columnstore.NewTableConfig(profile.SchemaDefinition()))
	require.NoError(t, err)

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	_, err = api.QueryRange(ctx, &pb.QueryRangeRequest{
		Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
		Start: timestamppb.New(timestamp.Time(0)),
		End:   timestamppb.New(timestamp.Time(9223372036854775807)),
	})
	require.ErrorIs(t, err, status.Error(
		codes.NotFound,
		"No data found for the query, try a different query or time range or no data has been written to be queried yet.",
	))
}

type Testing interface {
	require.TestingT
	Helper()
}

func MustReadAllGzip(t Testing, filename string) []byte {
	t.Helper()

	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()

	r, err := gzip.NewReader(f)
	require.NoError(t, err)
	content, err := io.ReadAll(r)
	require.NoError(t, err)
	return content
}

func MustDecompressGzip(t Testing, b []byte) []byte {
	t.Helper()

	r, err := gzip.NewReader(bytes.NewReader(b))
	require.NoError(t, err)
	content, err := io.ReadAll(r)
	require.NoError(t, err)
	return content
}

func TestColumnQueryAPIQueryRange(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	dir := "./testdata/many/"
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)
	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	for _, f := range files {
		fileContent, err := os.ReadFile(dir + f.Name())
		require.NoError(t, err)

		_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
			Series: []*profilestorepb.RawProfileSeries{{
				Labels: &profilestorepb.LabelSet{
					Labels: []*profilestorepb.Label{
						{
							Name:  "__name__",
							Value: "memory",
						},
						{
							Name:  "job",
							Value: "default",
						},
					},
				},
				Samples: []*profilestorepb.RawSample{{
					RawProfile: fileContent,
				}},
			}},
		})
		require.NoError(t, err)
	}

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	res, err := api.QueryRange(ctx, &pb.QueryRangeRequest{
		Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
		Start: timestamppb.New(timestamp.Time(0)),
		End:   timestamppb.New(timestamp.Time(9223372036854775807)),
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Series))
	require.Equal(t, 1, len(res.Series[0].Labelset.Labels))
	require.Equal(t, 10, len(res.Series[0].Samples))
}

func ptrToString(s string) *string {
	return &s
}

func TestColumnQueryAPIQuerySingle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)
	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)
	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	fileContent, err := os.ReadFile("testdata/alloc_objects.pb.gz")
	require.NoError(t, err)

	p := &pprofpb.Profile{}
	require.NoError(t, p.UnmarshalVT(MustDecompressGzip(t, fileContent)))

	_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: fileContent,
			}},
		}},
	})
	require.NoError(t, err)

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	ts := timestamppb.New(timestamp.Time(p.TimeNanos / time.Millisecond.Nanoseconds()))
	res, err := api.Query(ctx, &pb.QueryRequest{
		Options: &pb.QueryRequest_Single{
			Single: &pb.SingleProfile{
				Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
				Time:  ts,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, int32(33), res.Report.(*pb.QueryResponse_Flamegraph).Flamegraph.Height)

	res, err = api.Query(ctx, &pb.QueryRequest{
		ReportType: pb.QueryRequest_REPORT_TYPE_PPROF,
		Options: &pb.QueryRequest_Single{
			Single: &pb.SingleProfile{
				Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
				Time:  ts,
			},
		},
	})
	require.NoError(t, err)

	unfilteredRes, err := api.Query(ctx, &pb.QueryRequest{
		ReportType: pb.QueryRequest_REPORT_TYPE_TOP,
		Options: &pb.QueryRequest_Single{
			Single: &pb.SingleProfile{
				Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
				Time:  ts,
			},
		},
	})
	require.NoError(t, err)

	filteredRes, err := api.Query(ctx, &pb.QueryRequest{
		ReportType: pb.QueryRequest_REPORT_TYPE_TOP,
		Options: &pb.QueryRequest_Single{
			Single: &pb.SingleProfile{
				Query: `memory:alloc_objects:count:space:bytes{job="default", __name__="memory"}`,
				Time:  ts,
			},
		},
		FilterQuery: ptrToString("runtime"),
	})
	require.NoError(t, err)
	require.Less(t, len(filteredRes.Report.(*pb.QueryResponse_Top).Top.List), len(unfilteredRes.Report.(*pb.QueryResponse_Top).Top.List), "filtered result should be smaller than unfiltered result")

	testProf := &pprofpb.Profile{}
	err = testProf.UnmarshalVT(MustDecompressGzip(t, res.Report.(*pb.QueryResponse_Pprof).Pprof))
	require.NoError(t, err)
}

func TestColumnQueryAPIQueryFgprof(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	fileContent, err := os.ReadFile("testdata/fgprof.pb.gz")
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)

	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "fgprof",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: fileContent,
			}},
		}},
	})
	require.NoError(t, err)

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)

	res, err := api.QueryRange(ctx, &pb.QueryRangeRequest{
		Query: `fgprof:samples:count:wallclock:nanoseconds:delta`,
		Start: timestamppb.New(timestamp.Time(0)),
		End:   timestamppb.New(timestamp.Time(9223372036854775807)),
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Series))
	require.Equal(t, 1, len(res.Series[0].Labelset.Labels))
	require.Equal(t, 1, len(res.Series[0].Samples))
}

func TestColumnQueryAPIQueryCumulative(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)
	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	// Load CPU and memory profiles
	fileNames := []string{
		"testdata/alloc_objects.pb.gz",
		"testdata/profile1.pb.gz",
	}
	labelSets := []*profilestorepb.LabelSet{
		{
			Labels: []*profilestorepb.Label{
				{Name: "__name__", Value: "memory"},
				{Name: "job", Value: "default"},
			},
		},
		{
			Labels: []*profilestorepb.Label{
				{Name: "__name__", Value: "cpu"},
				{Name: "job", Value: "default"},
			},
		},
	}
	for i, fileName := range fileNames {
		fileContent, err := os.ReadFile(fileName)
		require.NoError(t, err)

		p := &pprofpb.Profile{}
		require.NoError(t, p.UnmarshalVT(MustDecompressGzip(t, fileContent)))

		_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
			Series: []*profilestorepb.RawProfileSeries{{
				Labels: labelSets[i],
				Samples: []*profilestorepb.RawSample{{
					RawProfile: fileContent,
				}},
			}},
		})
		require.NoError(t, err)
	}

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)

	// These have been extracted from the profiles above.
	queries := []struct {
		name      string
		query     string
		timeNanos int64
		// expected
		total    int64
		filtered int64
	}{{
		name:      "memory",
		query:     `memory:alloc_objects:count:space:bytes{job="default"}`,
		timeNanos: 1608199718549304626,
		total:     int64(310797348),
		filtered:  int64(0),
	}, {
		name:      "cpu",
		query:     `cpu:samples:count:cpu:nanoseconds:delta{job="default"}`,
		timeNanos: 1626013307085084416,
		total:     int64(48),
		filtered:  int64(0),
	}}

	// Check that the following report type return the same cumulative and filtered values.

	reportTypes := []pb.QueryRequest_ReportType{
		pb.QueryRequest_REPORT_TYPE_TOP,
		pb.QueryRequest_REPORT_TYPE_CALLGRAPH,
		pb.QueryRequest_REPORT_TYPE_FLAMEGRAPH_TABLE,
		pb.QueryRequest_REPORT_TYPE_FLAMEGRAPH_ARROW,
	}

	for _, query := range queries {
		for _, reportType := range reportTypes {
			t.Run(query.name+"-"+pb.QueryRequest_ReportType_name[int32(reportType)], func(t *testing.T) {
				res, err := api.Query(ctx, &pb.QueryRequest{
					ReportType: pb.QueryRequest_REPORT_TYPE_TOP,
					Options: &pb.QueryRequest_Single{
						Single: &pb.SingleProfile{
							Query: query.query,
							Time:  timestamppb.New(timestamp.Time(query.timeNanos / time.Millisecond.Nanoseconds())),
						},
					},
				})
				require.NoError(t, err)
				require.Equal(t, query.total, res.Total)
				require.Equal(t, query.filtered, res.Filtered)
			})
		}
	}
}

func MustCompressGzip(t Testing, p *pprofpb.Profile) []byte {
	t.Helper()

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	content, err := p.MarshalVT()
	require.NoError(t, err)
	_, err = w.Write(content)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestColumnQueryAPIQueryDiff(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	p := &pprofpb.Profile{
		StringTable: []string{
			"",
			"testFunc",
			"alloc_objects",
			"count",
			"space",
			"bytes",
		},
		Function: []*pprofpb.Function{{
			Id:   1,
			Name: 1,
		}},
		Location: []*pprofpb.Location{{
			Id:      1,
			Address: 0x1,
			Line: []*pprofpb.Line{{
				Line:       1,
				FunctionId: 1,
			}},
		}, {
			Id:      2,
			Address: 0x2,
			Line: []*pprofpb.Line{{
				Line:       2,
				FunctionId: 1,
			}},
		}},
		SampleType: []*pprofpb.ValueType{{
			Type: 2,
			Unit: 3,
		}},
		PeriodType: &pprofpb.ValueType{
			Type: 4,
			Unit: 5,
		},
		TimeNanos: 1000000,
		Sample: []*pprofpb.Sample{{
			Value:      []int64{1},
			LocationId: []uint64{1},
		}},
	}

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)

	r, err := normalizer.NormalizeWriteRawRequest(ctx, normalizer.New(), &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: MustCompressGzip(t, p),
			}},
		}},
	})
	require.NoError(t, err)
	require.NoError(t, ingester.Ingest(ctx, r))

	p.Sample = []*pprofpb.Sample{{
		Value:      []int64{2},
		LocationId: []uint64{2},
	}}
	p.TimeNanos = 2000000
	r, err = normalizer.NormalizeWriteRawRequest(ctx, normalizer.New(), &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: MustCompressGzip(t, p),
			}},
		}},
	})
	require.NoError(t, err)
	require.NoError(t, ingester.Ingest(ctx, r))

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)

	res, err := api.Query(ctx, &pb.QueryRequest{
		Mode: pb.QueryRequest_MODE_DIFF,
		Options: &pb.QueryRequest_Diff{
			Diff: &pb.DiffProfile{
				A: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(1)),
						},
					},
				},
				B: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(2)),
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	fg := res.Report.(*pb.QueryResponse_Flamegraph).Flamegraph
	require.Equal(t, int32(2), fg.Height)
	require.Equal(t, 1, len(fg.Root.Children))
	require.Equal(t, int64(2), fg.Root.Children[0].Cumulative)
	require.Equal(t, int64(1), fg.Root.Children[0].Diff)

	res, err = api.Query(ctx, &pb.QueryRequest{
		Mode:       pb.QueryRequest_MODE_DIFF,
		ReportType: *pb.QueryRequest_REPORT_TYPE_TOP.Enum(),
		Options: &pb.QueryRequest_Diff{
			Diff: &pb.DiffProfile{
				A: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(1)),
						},
					},
				},
				B: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(2)),
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	topList := res.Report.(*pb.QueryResponse_Top).Top.List
	require.Equal(t, 1, len(topList))
	require.Equal(t, int64(2), topList[0].Cumulative)
	require.Equal(t, int64(1), topList[0].Diff)

	res, err = api.Query(ctx, &pb.QueryRequest{
		Mode:       pb.QueryRequest_MODE_DIFF,
		ReportType: *pb.QueryRequest_REPORT_TYPE_PPROF.Enum(),
		Options: &pb.QueryRequest_Diff{
			Diff: &pb.DiffProfile{
				A: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(1)),
						},
					},
				},
				B: &pb.ProfileDiffSelection{
					Mode: pb.ProfileDiffSelection_MODE_SINGLE_UNSPECIFIED,
					Options: &pb.ProfileDiffSelection_Single{
						Single: &pb.SingleProfile{
							Query: `memory:alloc_objects:count:space:bytes{job="default"}`,
							Time:  timestamppb.New(timestamp.Time(2)),
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	testProf := &pprofpb.Profile{}
	err = testProf.UnmarshalVT(MustDecompressGzip(t, res.Report.(*pb.QueryResponse_Pprof).Pprof))
	require.NoError(t, err)
	require.Equal(t, 2, len(testProf.Sample))
	require.Equal(t, []int64{2}, testProf.Sample[0].Value)
	require.Equal(t, []int64{-1}, testProf.Sample[1].Value)
}

func TestColumnQueryAPITypes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	fileContent, err := os.ReadFile("testdata/alloc_space_delta.pb.gz")
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)

	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: fileContent,
			}},
		}},
	})
	require.NoError(t, err)

	require.NoError(t, table.EnsureCompaction())

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	res, err := api.ProfileTypes(ctx, &pb.ProfileTypesRequest{})
	require.NoError(t, err)

	/* res returned by profile type on arm machine did not have same ordering
	on `SampleType: "inuse_objects"` and `inuse_space`. Due to which test
	was quite flaky and failing. So instead of testing for exact structure of
	the proto message, comparing by proto size of the messages.
	*/
	require.Equal(t, proto.Size(&pb.ProfileTypesResponse{Types: []*pb.ProfileType{
		{Name: "memory", SampleType: "alloc_objects", SampleUnit: "count", PeriodType: "space", PeriodUnit: "bytes", Delta: true},
		{Name: "memory", SampleType: "alloc_space", SampleUnit: "bytes", PeriodType: "space", PeriodUnit: "bytes", Delta: true},
		{Name: "memory", SampleType: "inuse_objects", SampleUnit: "count", PeriodType: "space", PeriodUnit: "bytes", Delta: true},
		{Name: "memory", SampleType: "inuse_space", SampleUnit: "bytes", PeriodType: "space", PeriodUnit: "bytes", Delta: true},
	}}), proto.Size(res))
}

func TestColumnQueryAPILabelNames(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	fileContent, err := os.ReadFile("testdata/alloc_objects.pb.gz")
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)

	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: fileContent,
			}},
		}},
	})
	require.NoError(t, err)

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	res, err := api.Labels(ctx, &pb.LabelsRequest{})
	require.NoError(t, err)

	require.Equal(t, []string{
		"job",
	}, res.LabelNames)
}

func TestColumnQueryAPILabelValues(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := log.NewNopLogger()
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("")
	col, err := columnstore.New()
	require.NoError(t, err)
	colDB, err := col.DB(context.Background(), "parca")
	require.NoError(t, err)

	schema, err := profile.Schema()
	require.NoError(t, err)

	table, err := colDB.Table(
		"stacktraces",
		columnstore.NewTableConfig(profile.SchemaDefinition()),
	)
	require.NoError(t, err)

	fileContent, err := os.ReadFile("testdata/alloc_objects.pb.gz")
	require.NoError(t, err)

	ingester := ingester.NewIngester(
		logger,
		memory.DefaultAllocator,
		table,
		schema,
	)
	store := profilestore.NewProfileColumnStore(
		reg,
		logger,
		tracer,
		ingester,
		true,
	)

	_, err = store.WriteRaw(ctx, &profilestorepb.WriteRawRequest{
		Series: []*profilestorepb.RawProfileSeries{{
			Labels: &profilestorepb.LabelSet{
				Labels: []*profilestorepb.Label{
					{
						Name:  "__name__",
						Value: "memory",
					},
					{
						Name:  "job",
						Value: "default",
					},
				},
			},
			Samples: []*profilestorepb.RawSample{{
				RawProfile: fileContent,
			}},
		}},
	})
	require.NoError(t, err)

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	api := NewColumnQueryAPI(
		logger,
		tracer,
		getShareServerConn(t),
		parcacol.NewQuerier(
			logger,
			tracer,
			query.NewEngine(
				mem,
				colDB.TableProvider(),
			),
			"stacktraces",
			nil,
			mem,
		),
		mem,
		parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
		nil,
	)
	res, err := api.Values(ctx, &pb.ValuesRequest{
		LabelName: "job",
	})
	require.NoError(t, err)

	require.Equal(t, []string{
		"default",
	}, res.LabelValues)
}

func BenchmarkQuery(b *testing.B) {
	ctx := context.Background()
	tracer := noop.NewTracerProvider().Tracer("")

	fileContent, err := os.ReadFile("testdata/alloc_objects.pb.gz")
	require.NoError(b, err)

	p, err := pprofprofile.ParseData(fileContent)
	require.NoError(b, err)

	sp, err := PprofToSymbolizedProfile(profile.Meta{}, p, 0)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(b, 0)
	for i := 0; i < b.N; i++ {
		_, _ = RenderReport(
			ctx,
			tracer,
			sp,
			pb.QueryRequest_REPORT_TYPE_FLAMEGRAPH_ARROW,
			0,
			0,
			[]string{FlamegraphFieldFunctionName},
			NewTableConverterPool(),
			mem,
			parcacol.NewArrowToProfileConverter(tracer, kv.NewKeyMaker()),
			nil,
			"",
			false,
		)
	}
}

func PprofToSymbolizedProfile(meta profile.Meta, prof *pprofprofile.Profile, index int) (profile.Profile, error) {
	labelNameSet := make(map[string]struct{})
	for _, s := range prof.Sample {
		for k := range s.Label {
			labelNameSet[k] = struct{}{}
		}
	}
	labelNames := make([]string, 0, len(labelNameSet))
	for l := range labelNameSet {
		labelNames = append(labelNames, l)
	}

	w := profile.NewWriter(memory.DefaultAllocator, labelNames)
	defer w.RecordBuilder.Release()
	for i := range prof.Sample {
		if len(prof.Sample[i].Value) <= index {
			return profile.Profile{}, status.Errorf(codes.InvalidArgument, "failed to find samples for profile type")
		}

		w.Value.Append(prof.Sample[i].Value[index])
		w.ValuePerSecond.Append(0)
		w.Diff.Append(0)
		w.DiffPerSecond.Append(0)

		for labelName, labelBuilder := range w.LabelBuildersMap {
			if prof.Sample[i].Label == nil {
				labelBuilder.AppendNull()
				continue
			}

			if labelValues, ok := prof.Sample[i].Label[labelName]; ok && len(labelValues) > 0 {
				labelBuilder.Append([]byte(labelValues[0]))
			} else {
				labelBuilder.AppendNull()
			}
		}

		w.LocationsList.Append(len(prof.Sample[i].Location) > 0)
		if len(prof.Sample[i].Location) > 0 {
			for _, loc := range prof.Sample[i].Location {
				w.Locations.Append(true)
				w.Addresses.Append(loc.Address)

				if loc.Mapping != nil {
					w.MappingStart.Append(loc.Mapping.Start)
					w.MappingLimit.Append(loc.Mapping.Limit)
					w.MappingOffset.Append(loc.Mapping.Offset)
					w.MappingFile.Append([]byte(loc.Mapping.File))
					w.MappingBuildID.Append([]byte(loc.Mapping.BuildID))
				} else {
					w.MappingStart.AppendNull()
					w.MappingLimit.AppendNull()
					w.MappingOffset.AppendNull()
					w.MappingFile.AppendNull()
					w.MappingBuildID.AppendNull()
				}

				w.Lines.Append(len(loc.Line) > 0)
				if len(loc.Line) > 0 {
					for _, line := range loc.Line {
						w.Line.Append(true)
						w.LineNumber.Append(line.Line)
						if line.Function != nil {
							w.FunctionName.Append([]byte(line.Function.Name))
							w.FunctionSystemName.Append([]byte(line.Function.SystemName))
							w.FunctionFilename.Append([]byte(line.Function.Filename))
							w.FunctionStartLine.Append(line.Function.StartLine)
						} else {
							w.FunctionName.AppendNull()
							w.FunctionSystemName.AppendNull()
							w.FunctionFilename.AppendNull()
							w.FunctionStartLine.AppendNull()
						}
					}
				}
			}
		}
	}

	return profile.Profile{
		Meta:    meta,
		Samples: []arrow.Record{w.RecordBuilder.NewRecord()},
	}, nil
}

func OldProfileToArrowProfile(p profile.OldProfile) (profile.Profile, error) {
	labelNameSet := make(map[string]struct{})
	for _, s := range p.Samples {
		for k := range s.Label {
			labelNameSet[k] = struct{}{}
		}
	}
	labelNames := make([]string, 0, len(labelNameSet))
	for l := range labelNameSet {
		labelNames = append(labelNames, l)
	}

	w := profile.NewWriter(memory.DefaultAllocator, labelNames)
	defer w.RecordBuilder.Release()
	for i := range p.Samples {
		w.Value.Append(p.Samples[i].Value)
		w.ValuePerSecond.Append(float64(p.Samples[i].Value))
		w.Diff.Append(p.Samples[i].DiffValue)
		w.DiffPerSecond.Append(float64(p.Samples[i].DiffValue))

		for labelName, labelBuilder := range w.LabelBuildersMap {
			if p.Samples[i].Label == nil {
				labelBuilder.AppendNull()
				continue
			}

			if labelValue, ok := p.Samples[i].Label[labelName]; ok {
				labelBuilder.Append([]byte(labelValue))
			} else {
				labelBuilder.AppendNull()
			}
		}

		w.LocationsList.Append(len(p.Samples[i].Locations) > 0)
		if len(p.Samples[i].Locations) > 0 {
			for _, loc := range p.Samples[i].Locations {
				w.Locations.Append(true)
				w.Addresses.Append(loc.Address)

				if loc.Mapping != nil {
					w.MappingStart.Append(loc.Mapping.Start)
					w.MappingLimit.Append(loc.Mapping.Limit)
					w.MappingOffset.Append(loc.Mapping.Offset)
					w.MappingFile.Append([]byte(loc.Mapping.File))
					w.MappingBuildID.Append([]byte(loc.Mapping.BuildId))
				} else {
					w.MappingStart.AppendNull()
					w.MappingLimit.AppendNull()
					w.MappingOffset.AppendNull()
					w.MappingFile.AppendNull()
					w.MappingBuildID.AppendNull()
				}

				w.Lines.Append(len(loc.Lines) > 0)
				if len(loc.Lines) > 0 {
					for _, line := range loc.Lines {
						if line.Function != nil {
							w.Line.Append(true)
							w.LineNumber.Append(line.Line)
							w.FunctionName.Append([]byte(line.Function.Name))
							w.FunctionSystemName.Append([]byte(line.Function.SystemName))
							w.FunctionFilename.Append([]byte(line.Function.Filename))
							w.FunctionStartLine.Append(line.Function.StartLine)
						} else {
							w.Line.AppendNull()
						}
					}
				}
			}
		}
	}

	return profile.Profile{
		Meta:    p.Meta,
		Samples: []arrow.Record{w.RecordBuilder.NewRecord()},
	}, nil
}

func TestFilterData(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	w := profile.NewWriter(mem, nil)
	defer w.Release()

	w.LocationsList.Append(true)
	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("test"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(1)
	w.FunctionName.Append([]byte("test"))
	w.FunctionSystemName.Append([]byte("test"))
	w.FunctionFilename.Append([]byte("test"))
	w.FunctionStartLine.Append(1)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("libpython3.11.so.1.0"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(1)
	w.FunctionName.Append([]byte("test1"))
	w.FunctionSystemName.Append([]byte("test"))
	w.FunctionFilename.Append([]byte("test"))
	w.FunctionStartLine.Append(1)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("test"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(1)
	w.FunctionName.Append([]byte("test1"))
	w.FunctionSystemName.Append([]byte("test"))
	w.FunctionFilename.Append([]byte("test"))
	w.FunctionStartLine.Append(1)
	w.Value.Append(1)
	w.ValuePerSecond.Append(1)
	w.Diff.Append(0)
	w.DiffPerSecond.Append(0)

	originalRecord := w.RecordBuilder.NewRecord()
	recs, _, err := FilterProfileData(
		context.Background(),
		noop.NewTracerProvider().Tracer(""),
		mem,
		[]arrow.Record{originalRecord},
		"",
		&pb.RuntimeFilter{
			ShowPython: false,
		},
	)
	require.NoError(t, err)
	defer func() {
		for _, r := range recs {
			r.Release()
		}
	}()
	r := profile.NewRecordReader(recs[0])
	valid := 0
	for i := 0; i < r.Location.Len(); i++ {
		if r.Location.IsValid(i) {
			valid++
		}
	}
	require.Equal(t, 2, valid)
	require.Equal(t, "test", string(r.LineFunctionNameDict.Value(int(r.LineFunctionNameIndices.Value(0)))))
	require.Equal(t, "test1", string(r.LineFunctionNameDict.Value(int(r.LineFunctionNameIndices.Value(1)))))
}

func TestFilterUnsymbolized(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	w := profile.NewWriter(mem, nil)
	defer w.Release()

	w.LocationsList.Append(true)
	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("test"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(false)
	w.Value.Append(1)
	w.ValuePerSecond.Append(1)
	w.Diff.Append(0)
	w.DiffPerSecond.Append(0)

	originalRecord := w.RecordBuilder.NewRecord()
	recs, _, err := FilterProfileData(
		context.Background(),
		noop.NewTracerProvider().Tracer(""),
		mem,
		[]arrow.Record{originalRecord},
		"",
		&pb.RuntimeFilter{
			ShowPython: false,
		},
	)
	require.NoError(t, err)
	require.Len(t, recs, 1)
	defer func() {
		for _, r := range recs {
			r.Release()
		}
	}()
	r := profile.NewRecordReader(recs[0])
	valid := 0
	for i := 0; i < r.Location.Len(); i++ {
		if r.Location.IsValid(i) {
			valid++
		}
	}
	require.Equal(t, 1, valid)
}

func TestFilterDataWithPath(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	w := profile.NewWriter(mem, nil)
	defer w.Release()

	w.LocationsList.Append(true)
	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("libc.so.6"))
	w.MappingBuildID.Append([]byte(""))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(1)
	w.FunctionName.Append([]byte("__libc_start_main"))
	w.FunctionSystemName.Append([]byte("__libc_start_main"))
	w.FunctionFilename.Append([]byte(""))
	w.FunctionStartLine.Append(1)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("/usr/lib/libpython3.11.so.1.0"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(0)
	w.FunctionName.Append([]byte("test1"))
	w.FunctionSystemName.Append([]byte("test1"))
	w.FunctionFilename.Append([]byte(""))
	w.FunctionStartLine.Append(0)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("interpreter"))
	w.MappingBuildID.Append([]byte(""))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(0)
	w.FunctionName.Append([]byte("test"))
	w.FunctionSystemName.Append([]byte("test"))
	w.FunctionFilename.Append([]byte("test.py"))
	w.FunctionStartLine.Append(0)
	w.Value.Append(1)
	w.ValuePerSecond.Append(1)
	w.Diff.Append(0)
	w.DiffPerSecond.Append(0)

	originalRecord := w.RecordBuilder.NewRecord()
	recs, _, err := FilterProfileData(
		context.Background(),
		noop.NewTracerProvider().Tracer(""),
		mem,
		[]arrow.Record{originalRecord},
		"",
		nil,
	)
	require.NoError(t, err)
	defer func() {
		for _, r := range recs {
			r.Release()
		}
	}()
	r := profile.NewRecordReader(recs[0])
	valid := 0
	for i := 0; i < r.Location.Len(); i++ {
		if r.Location.IsValid(i) {
			valid++
		}
	}
	require.Equal(t, 2, valid)
	require.Equal(t, "__libc_start_main", string(r.LineFunctionNameDict.Value(int(r.LineFunctionNameIndices.Value(0)))))
	require.Equal(t, "test", string(r.LineFunctionNameDict.Value(int(r.LineFunctionNameIndices.Value(2)))))
}

func TestFilterDataInterpretedOnly(t *testing.T) {
	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	w := profile.NewWriter(mem, nil)
	defer w.Release()

	w.LocationsList.Append(true)
	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("libc.so.6"))
	w.MappingBuildID.Append([]byte(""))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(1)
	w.FunctionName.Append([]byte("__libc_start_main"))
	w.FunctionSystemName.Append([]byte("__libc_start_main"))
	w.FunctionFilename.Append([]byte(""))
	w.FunctionStartLine.Append(1)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("/usr/lib/libpython3.11.so.1.0"))
	w.MappingBuildID.Append([]byte("test"))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(0)
	w.FunctionName.Append([]byte("test1"))
	w.FunctionSystemName.Append([]byte("test1"))
	w.FunctionFilename.Append([]byte(""))
	w.FunctionStartLine.Append(0)

	w.Locations.Append(true)
	w.Addresses.Append(0x1234)
	w.MappingStart.Append(0x1000)
	w.MappingLimit.Append(0x2000)
	w.MappingOffset.Append(0x0)
	w.MappingFile.Append([]byte("interpreter"))
	w.MappingBuildID.Append([]byte(""))
	w.Lines.Append(true)
	w.Line.Append(true)
	w.LineNumber.Append(0)
	w.FunctionName.Append([]byte("test"))
	w.FunctionSystemName.Append([]byte("test"))
	w.FunctionFilename.Append([]byte("test.py"))
	w.FunctionStartLine.Append(0)
	w.Value.Append(1)
	w.ValuePerSecond.Append(1)
	w.Diff.Append(0)
	w.DiffPerSecond.Append(0)

	originalRecord := w.RecordBuilder.NewRecord()
	recs, _, err := FilterProfileData(
		context.Background(),
		noop.NewTracerProvider().Tracer(""),
		mem,
		[]arrow.Record{originalRecord},
		"",
		&pb.RuntimeFilter{
			ShowInterpretedOnly: true,
		},
	)
	require.NoError(t, err)
	defer func() {
		for _, r := range recs {
			r.Release()
		}
	}()
	r := profile.NewRecordReader(recs[0])
	valid := 0
	for i := 0; i < r.Location.Len(); i++ {
		if r.Location.IsValid(i) {
			valid++
		}
	}
	require.Equal(t, 1, valid)
	require.Equal(t, "test", string(r.LineFunctionNameDict.Value(int(r.LineFunctionNameIndices.Value(2)))))
}

func BenchmarkFilterData(t *testing.B) {
	mem := memory.NewCheckedAllocator(memory.DefaultAllocator)
	defer mem.AssertSize(t, 0)
	w := profile.NewWriter(mem, nil)
	defer w.Release()

	for i := 0; i < 10000; i++ {
		w.LocationsList.Append(true)
		w.Locations.Append(true)
		w.Addresses.Append(0x1234)
		w.MappingStart.Append(0x1000)
		w.MappingLimit.Append(0x2000)
		w.MappingOffset.Append(0x0)
		w.MappingFile.Append([]byte("test"))
		w.MappingBuildID.Append([]byte("test"))
		w.Lines.Append(true)
		w.Line.Append(true)
		w.LineNumber.Append(1)
		w.FunctionName.Append([]byte("test"))
		w.FunctionSystemName.Append([]byte("test"))
		w.FunctionFilename.Append([]byte("test"))
		w.FunctionStartLine.Append(1)

		w.Locations.Append(true)
		w.Addresses.Append(0x1234)
		w.MappingStart.Append(0x1000)
		w.MappingLimit.Append(0x2000)
		w.MappingOffset.Append(0x0)
		w.MappingFile.Append([]byte("libpython3.11.so.1.0"))
		w.MappingBuildID.Append([]byte("test"))
		w.Lines.Append(true)
		w.Line.Append(true)
		w.LineNumber.Append(1)
		w.FunctionName.Append([]byte("test1"))
		w.FunctionSystemName.Append([]byte("test"))
		w.FunctionFilename.Append([]byte("test"))
		w.FunctionStartLine.Append(1)

		w.Locations.Append(true)
		w.Addresses.Append(0x1234)
		w.MappingStart.Append(0x1000)
		w.MappingLimit.Append(0x2000)
		w.MappingOffset.Append(0x0)
		w.MappingFile.Append([]byte("test"))
		w.MappingBuildID.Append([]byte("test"))
		w.Lines.Append(true)
		w.Line.Append(true)
		w.LineNumber.Append(1)
		w.FunctionName.Append([]byte("test1"))
		w.FunctionSystemName.Append([]byte("test"))
		w.FunctionFilename.Append([]byte("test"))
		w.FunctionStartLine.Append(1)
		w.Value.Append(1)
		w.ValuePerSecond.Append(1)
		w.Diff.Append(0)
		w.DiffPerSecond.Append(0)
	}

	originalRecord := w.RecordBuilder.NewRecord()
	defer originalRecord.Release()
	for i := 0; i < t.N; i++ {
		originalRecord.Retain() // retain each time since FilterProfileData will release it
		recs, _, err := FilterProfileData(
			context.Background(),
			noop.NewTracerProvider().Tracer(""),
			mem,
			[]arrow.Record{originalRecord},
			"",
			&pb.RuntimeFilter{
				ShowPython: false,
			},
		)
		require.NoError(t, err)
		for _, r := range recs {
			r.Release()
		}
	}
}
