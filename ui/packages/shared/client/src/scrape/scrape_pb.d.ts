// package: parca.scrape
// file: scrape/scrape.proto

import * as jspb from "google-protobuf";
import * as google_api_annotations_pb from "../google/api/annotations_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as profilestore_profilestore_pb from "../profilestore/profilestore_pb";

export class TargetsRequest extends jspb.Message {
  getState(): TargetsRequest.StateMap[keyof TargetsRequest.StateMap];
  setState(value: TargetsRequest.StateMap[keyof TargetsRequest.StateMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsRequest): TargetsRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TargetsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsRequest;
  static deserializeBinaryFromReader(message: TargetsRequest, reader: jspb.BinaryReader): TargetsRequest;
}

export namespace TargetsRequest {
  export type AsObject = {
    state: TargetsRequest.StateMap[keyof TargetsRequest.StateMap],
  }

  export interface StateMap {
    ANY: 0;
    ACTIVE: 1;
    DROPPED: 2;
  }

  export const State: StateMap;
}

export class TargetsResponse extends jspb.Message {
  getTargetsMap(): jspb.Map<string, Targets>;
  clearTargetsMap(): void;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TargetsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TargetsResponse): TargetsResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TargetsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TargetsResponse;
  static deserializeBinaryFromReader(message: TargetsResponse, reader: jspb.BinaryReader): TargetsResponse;
}

export namespace TargetsResponse {
  export type AsObject = {
    targetsMap: Array<[string, Targets.AsObject]>,
  }
}

export class Targets extends jspb.Message {
  clearTargetsList(): void;
  getTargetsList(): Array<Target>;
  setTargetsList(value: Array<Target>): void;
  addTargets(value?: Target, index?: number): Target;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Targets.AsObject;
  static toObject(includeInstance: boolean, msg: Targets): Targets.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Targets, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Targets;
  static deserializeBinaryFromReader(message: Targets, reader: jspb.BinaryReader): Targets;
}

export namespace Targets {
  export type AsObject = {
    targetsList: Array<Target.AsObject>,
  }
}

export class Target extends jspb.Message {
  hasDiscoveredlabels(): boolean;
  clearDiscoveredlabels(): void;
  getDiscoveredlabels(): profilestore_profilestore_pb.LabelSet | undefined;
  setDiscoveredlabels(value?: profilestore_profilestore_pb.LabelSet): void;

  hasLabels(): boolean;
  clearLabels(): void;
  getLabels(): profilestore_profilestore_pb.LabelSet | undefined;
  setLabels(value?: profilestore_profilestore_pb.LabelSet): void;

  getLasterror(): string;
  setLasterror(value: string): void;

  hasLastscrape(): boolean;
  clearLastscrape(): void;
  getLastscrape(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastscrape(value?: google_protobuf_timestamp_pb.Timestamp): void;

  hasLastscrapeduration(): boolean;
  clearLastscrapeduration(): void;
  getLastscrapeduration(): google_protobuf_duration_pb.Duration | undefined;
  setLastscrapeduration(value?: google_protobuf_duration_pb.Duration): void;

  getUrl(): string;
  setUrl(value: string): void;

  getHealth(): Target.HealthMap[keyof Target.HealthMap];
  setHealth(value: Target.HealthMap[keyof Target.HealthMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Target.AsObject;
  static toObject(includeInstance: boolean, msg: Target): Target.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Target, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Target;
  static deserializeBinaryFromReader(message: Target, reader: jspb.BinaryReader): Target;
}

export namespace Target {
  export type AsObject = {
    discoveredlabels?: profilestore_profilestore_pb.LabelSet.AsObject,
    labels?: profilestore_profilestore_pb.LabelSet.AsObject,
    lasterror: string,
    lastscrape?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastscrapeduration?: google_protobuf_duration_pb.Duration.AsObject,
    url: string,
    health: Target.HealthMap[keyof Target.HealthMap],
  }

  export interface HealthMap {
    UNKNOWN: 0;
    GOOD: 1;
    BAD: 2;
  }

  export const Health: HealthMap;
}

