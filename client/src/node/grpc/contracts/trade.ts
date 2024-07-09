/* eslint-disable */
import {
  ChannelCredentials,
  Client,
  ClientReadableStream,
  handleServerStreamingCall,
  makeGenericClientConstructor,
  Metadata,
} from "@grpc/grpc-js";
import type {
  CallOptions,
  ClientOptions,
  ClientUnaryCall,
  handleUnaryCall,
  ServiceError,
  UntypedServiceImplementation,
} from "@grpc/grpc-js";
import _m0 from "protobufjs/minimal";
import { Struct } from "./google/protobuf/struct";

export const protobufPackage = "trade";

export interface SubscribeStrategiesEventsRequest {
  Strategy: string;
}

export interface StrategyEvent {
  ID: string;
  Kind: string;
  Value: { [key: string]: any } | undefined;
}

export interface StartRequest {
  Strategy: string;
  InstrumentId: string;
}

export interface StartResponse {
  Ok: boolean;
  Error: string;
}

export interface StopRequest {
  Strategy: string;
  InstrumentId: string;
}

export interface StopResponse {
  Ok: boolean;
  Error: string;
}

export interface IsStartedRequest {
  Strategy: string;
  InstrumentId: string;
}

export interface IsStartedResponse {
  Ok: boolean;
  Error: string;
}

export interface ChangeConfigRequest {
  Strategy: string;
  InstrumentId: string;
  Config: { [key: string]: any } | undefined;
}

export interface ChangeConfigResponse {
  Ok: boolean;
  Error: string;
}

export interface GetConfigRequest {
  Strategy: string;
  InstrumentId: string;
}

export interface GetConfigResponse {
  Config: { [key: string]: any } | undefined;
}

function createBaseSubscribeStrategiesEventsRequest(): SubscribeStrategiesEventsRequest {
  return { Strategy: "" };
}

export const SubscribeStrategiesEventsRequest = {
  encode(message: SubscribeStrategiesEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeStrategiesEventsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubscribeStrategiesEventsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubscribeStrategiesEventsRequest {
    return { Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "" };
  },

  toJSON(message: SubscribeStrategiesEventsRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubscribeStrategiesEventsRequest>, I>>(
    base?: I,
  ): SubscribeStrategiesEventsRequest {
    return SubscribeStrategiesEventsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubscribeStrategiesEventsRequest>, I>>(
    object: I,
  ): SubscribeStrategiesEventsRequest {
    const message = createBaseSubscribeStrategiesEventsRequest();
    message.Strategy = object.Strategy ?? "";
    return message;
  },
};

function createBaseStrategyEvent(): StrategyEvent {
  return { ID: "", Kind: "", Value: undefined };
}

export const StrategyEvent = {
  encode(message: StrategyEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ID !== "") {
      writer.uint32(10).string(message.ID);
    }
    if (message.Kind !== "") {
      writer.uint32(18).string(message.Kind);
    }
    if (message.Value !== undefined) {
      Struct.encode(Struct.wrap(message.Value), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StrategyEvent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStrategyEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.ID = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Kind = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Value = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StrategyEvent {
    return {
      ID: isSet(object.ID) ? globalThis.String(object.ID) : "",
      Kind: isSet(object.Kind) ? globalThis.String(object.Kind) : "",
      Value: isObject(object.Value) ? object.Value : undefined,
    };
  },

  toJSON(message: StrategyEvent): unknown {
    const obj: any = {};
    if (message.ID !== "") {
      obj.ID = message.ID;
    }
    if (message.Kind !== "") {
      obj.Kind = message.Kind;
    }
    if (message.Value !== undefined) {
      obj.Value = message.Value;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StrategyEvent>, I>>(base?: I): StrategyEvent {
    return StrategyEvent.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StrategyEvent>, I>>(object: I): StrategyEvent {
    const message = createBaseStrategyEvent();
    message.ID = object.ID ?? "";
    message.Kind = object.Kind ?? "";
    message.Value = object.Value ?? undefined;
    return message;
  },
};

function createBaseStartRequest(): StartRequest {
  return { Strategy: "", InstrumentId: "" };
}

export const StartRequest = {
  encode(message: StartRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    if (message.InstrumentId !== "") {
      writer.uint32(18).string(message.InstrumentId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.InstrumentId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartRequest {
    return {
      Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "",
      InstrumentId: isSet(object.InstrumentId) ? globalThis.String(object.InstrumentId) : "",
    };
  },

  toJSON(message: StartRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    if (message.InstrumentId !== "") {
      obj.InstrumentId = message.InstrumentId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StartRequest>, I>>(base?: I): StartRequest {
    return StartRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StartRequest>, I>>(object: I): StartRequest {
    const message = createBaseStartRequest();
    message.Strategy = object.Strategy ?? "";
    message.InstrumentId = object.InstrumentId ?? "";
    return message;
  },
};

function createBaseStartResponse(): StartResponse {
  return { Ok: false, Error: "" };
}

export const StartResponse = {
  encode(message: StartResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Ok === true) {
      writer.uint32(8).bool(message.Ok);
    }
    if (message.Error !== "") {
      writer.uint32(18).string(message.Error);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.Ok = reader.bool();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Error = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartResponse {
    return {
      Ok: isSet(object.Ok) ? globalThis.Boolean(object.Ok) : false,
      Error: isSet(object.Error) ? globalThis.String(object.Error) : "",
    };
  },

  toJSON(message: StartResponse): unknown {
    const obj: any = {};
    if (message.Ok === true) {
      obj.Ok = message.Ok;
    }
    if (message.Error !== "") {
      obj.Error = message.Error;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StartResponse>, I>>(base?: I): StartResponse {
    return StartResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StartResponse>, I>>(object: I): StartResponse {
    const message = createBaseStartResponse();
    message.Ok = object.Ok ?? false;
    message.Error = object.Error ?? "";
    return message;
  },
};

function createBaseStopRequest(): StopRequest {
  return { Strategy: "", InstrumentId: "" };
}

export const StopRequest = {
  encode(message: StopRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    if (message.InstrumentId !== "") {
      writer.uint32(18).string(message.InstrumentId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StopRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStopRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.InstrumentId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StopRequest {
    return {
      Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "",
      InstrumentId: isSet(object.InstrumentId) ? globalThis.String(object.InstrumentId) : "",
    };
  },

  toJSON(message: StopRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    if (message.InstrumentId !== "") {
      obj.InstrumentId = message.InstrumentId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StopRequest>, I>>(base?: I): StopRequest {
    return StopRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StopRequest>, I>>(object: I): StopRequest {
    const message = createBaseStopRequest();
    message.Strategy = object.Strategy ?? "";
    message.InstrumentId = object.InstrumentId ?? "";
    return message;
  },
};

function createBaseStopResponse(): StopResponse {
  return { Ok: false, Error: "" };
}

export const StopResponse = {
  encode(message: StopResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Ok === true) {
      writer.uint32(8).bool(message.Ok);
    }
    if (message.Error !== "") {
      writer.uint32(18).string(message.Error);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StopResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStopResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.Ok = reader.bool();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Error = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StopResponse {
    return {
      Ok: isSet(object.Ok) ? globalThis.Boolean(object.Ok) : false,
      Error: isSet(object.Error) ? globalThis.String(object.Error) : "",
    };
  },

  toJSON(message: StopResponse): unknown {
    const obj: any = {};
    if (message.Ok === true) {
      obj.Ok = message.Ok;
    }
    if (message.Error !== "") {
      obj.Error = message.Error;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<StopResponse>, I>>(base?: I): StopResponse {
    return StopResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<StopResponse>, I>>(object: I): StopResponse {
    const message = createBaseStopResponse();
    message.Ok = object.Ok ?? false;
    message.Error = object.Error ?? "";
    return message;
  },
};

function createBaseIsStartedRequest(): IsStartedRequest {
  return { Strategy: "", InstrumentId: "" };
}

export const IsStartedRequest = {
  encode(message: IsStartedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    if (message.InstrumentId !== "") {
      writer.uint32(18).string(message.InstrumentId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsStartedRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsStartedRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.InstrumentId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IsStartedRequest {
    return {
      Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "",
      InstrumentId: isSet(object.InstrumentId) ? globalThis.String(object.InstrumentId) : "",
    };
  },

  toJSON(message: IsStartedRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    if (message.InstrumentId !== "") {
      obj.InstrumentId = message.InstrumentId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IsStartedRequest>, I>>(base?: I): IsStartedRequest {
    return IsStartedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IsStartedRequest>, I>>(object: I): IsStartedRequest {
    const message = createBaseIsStartedRequest();
    message.Strategy = object.Strategy ?? "";
    message.InstrumentId = object.InstrumentId ?? "";
    return message;
  },
};

function createBaseIsStartedResponse(): IsStartedResponse {
  return { Ok: false, Error: "" };
}

export const IsStartedResponse = {
  encode(message: IsStartedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Ok === true) {
      writer.uint32(8).bool(message.Ok);
    }
    if (message.Error !== "") {
      writer.uint32(18).string(message.Error);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsStartedResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsStartedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.Ok = reader.bool();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Error = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IsStartedResponse {
    return {
      Ok: isSet(object.Ok) ? globalThis.Boolean(object.Ok) : false,
      Error: isSet(object.Error) ? globalThis.String(object.Error) : "",
    };
  },

  toJSON(message: IsStartedResponse): unknown {
    const obj: any = {};
    if (message.Ok === true) {
      obj.Ok = message.Ok;
    }
    if (message.Error !== "") {
      obj.Error = message.Error;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IsStartedResponse>, I>>(base?: I): IsStartedResponse {
    return IsStartedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IsStartedResponse>, I>>(object: I): IsStartedResponse {
    const message = createBaseIsStartedResponse();
    message.Ok = object.Ok ?? false;
    message.Error = object.Error ?? "";
    return message;
  },
};

function createBaseChangeConfigRequest(): ChangeConfigRequest {
  return { Strategy: "", InstrumentId: "", Config: undefined };
}

export const ChangeConfigRequest = {
  encode(message: ChangeConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    if (message.InstrumentId !== "") {
      writer.uint32(18).string(message.InstrumentId);
    }
    if (message.Config !== undefined) {
      Struct.encode(Struct.wrap(message.Config), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChangeConfigRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChangeConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.InstrumentId = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Config = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChangeConfigRequest {
    return {
      Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "",
      InstrumentId: isSet(object.InstrumentId) ? globalThis.String(object.InstrumentId) : "",
      Config: isObject(object.Config) ? object.Config : undefined,
    };
  },

  toJSON(message: ChangeConfigRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    if (message.InstrumentId !== "") {
      obj.InstrumentId = message.InstrumentId;
    }
    if (message.Config !== undefined) {
      obj.Config = message.Config;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChangeConfigRequest>, I>>(base?: I): ChangeConfigRequest {
    return ChangeConfigRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChangeConfigRequest>, I>>(object: I): ChangeConfigRequest {
    const message = createBaseChangeConfigRequest();
    message.Strategy = object.Strategy ?? "";
    message.InstrumentId = object.InstrumentId ?? "";
    message.Config = object.Config ?? undefined;
    return message;
  },
};

function createBaseChangeConfigResponse(): ChangeConfigResponse {
  return { Ok: false, Error: "" };
}

export const ChangeConfigResponse = {
  encode(message: ChangeConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Ok === true) {
      writer.uint32(8).bool(message.Ok);
    }
    if (message.Error !== "") {
      writer.uint32(18).string(message.Error);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ChangeConfigResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChangeConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.Ok = reader.bool();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Error = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ChangeConfigResponse {
    return {
      Ok: isSet(object.Ok) ? globalThis.Boolean(object.Ok) : false,
      Error: isSet(object.Error) ? globalThis.String(object.Error) : "",
    };
  },

  toJSON(message: ChangeConfigResponse): unknown {
    const obj: any = {};
    if (message.Ok === true) {
      obj.Ok = message.Ok;
    }
    if (message.Error !== "") {
      obj.Error = message.Error;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ChangeConfigResponse>, I>>(base?: I): ChangeConfigResponse {
    return ChangeConfigResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ChangeConfigResponse>, I>>(object: I): ChangeConfigResponse {
    const message = createBaseChangeConfigResponse();
    message.Ok = object.Ok ?? false;
    message.Error = object.Error ?? "";
    return message;
  },
};

function createBaseGetConfigRequest(): GetConfigRequest {
  return { Strategy: "", InstrumentId: "" };
}

export const GetConfigRequest = {
  encode(message: GetConfigRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Strategy !== "") {
      writer.uint32(10).string(message.Strategy);
    }
    if (message.InstrumentId !== "") {
      writer.uint32(18).string(message.InstrumentId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetConfigRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Strategy = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.InstrumentId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetConfigRequest {
    return {
      Strategy: isSet(object.Strategy) ? globalThis.String(object.Strategy) : "",
      InstrumentId: isSet(object.InstrumentId) ? globalThis.String(object.InstrumentId) : "",
    };
  },

  toJSON(message: GetConfigRequest): unknown {
    const obj: any = {};
    if (message.Strategy !== "") {
      obj.Strategy = message.Strategy;
    }
    if (message.InstrumentId !== "") {
      obj.InstrumentId = message.InstrumentId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetConfigRequest>, I>>(base?: I): GetConfigRequest {
    return GetConfigRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetConfigRequest>, I>>(object: I): GetConfigRequest {
    const message = createBaseGetConfigRequest();
    message.Strategy = object.Strategy ?? "";
    message.InstrumentId = object.InstrumentId ?? "";
    return message;
  },
};

function createBaseGetConfigResponse(): GetConfigResponse {
  return { Config: undefined };
}

export const GetConfigResponse = {
  encode(message: GetConfigResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Config !== undefined) {
      Struct.encode(Struct.wrap(message.Config), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetConfigResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetConfigResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Config = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetConfigResponse {
    return { Config: isObject(object.Config) ? object.Config : undefined };
  },

  toJSON(message: GetConfigResponse): unknown {
    const obj: any = {};
    if (message.Config !== undefined) {
      obj.Config = message.Config;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetConfigResponse>, I>>(base?: I): GetConfigResponse {
    return GetConfigResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetConfigResponse>, I>>(object: I): GetConfigResponse {
    const message = createBaseGetConfigResponse();
    message.Config = object.Config ?? undefined;
    return message;
  },
};

export type TradeService = typeof TradeService;
export const TradeService = {
  start: {
    path: "/trade.Trade/Start",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: StartRequest) => Buffer.from(StartRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => StartRequest.decode(value),
    responseSerialize: (value: StartResponse) => Buffer.from(StartResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => StartResponse.decode(value),
  },
  stop: {
    path: "/trade.Trade/Stop",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: StopRequest) => Buffer.from(StopRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => StopRequest.decode(value),
    responseSerialize: (value: StopResponse) => Buffer.from(StopResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => StopResponse.decode(value),
  },
  isStarted: {
    path: "/trade.Trade/IsStarted",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: StartRequest) => Buffer.from(StartRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => StartRequest.decode(value),
    responseSerialize: (value: StartResponse) => Buffer.from(StartResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => StartResponse.decode(value),
  },
  changeConfig: {
    path: "/trade.Trade/ChangeConfig",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: ChangeConfigRequest) => Buffer.from(ChangeConfigRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => ChangeConfigRequest.decode(value),
    responseSerialize: (value: ChangeConfigResponse) => Buffer.from(ChangeConfigResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => ChangeConfigResponse.decode(value),
  },
  getConfig: {
    path: "/trade.Trade/GetConfig",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: GetConfigRequest) => Buffer.from(GetConfigRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => GetConfigRequest.decode(value),
    responseSerialize: (value: GetConfigResponse) => Buffer.from(GetConfigResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => GetConfigResponse.decode(value),
  },
  subscribeStrategiesEvents: {
    path: "/trade.Trade/SubscribeStrategiesEvents",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: SubscribeStrategiesEventsRequest) =>
      Buffer.from(SubscribeStrategiesEventsRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SubscribeStrategiesEventsRequest.decode(value),
    responseSerialize: (value: StrategyEvent) => Buffer.from(StrategyEvent.encode(value).finish()),
    responseDeserialize: (value: Buffer) => StrategyEvent.decode(value),
  },
} as const;

export interface TradeServer extends UntypedServiceImplementation {
  start: handleUnaryCall<StartRequest, StartResponse>;
  stop: handleUnaryCall<StopRequest, StopResponse>;
  isStarted: handleUnaryCall<StartRequest, StartResponse>;
  changeConfig: handleUnaryCall<ChangeConfigRequest, ChangeConfigResponse>;
  getConfig: handleUnaryCall<GetConfigRequest, GetConfigResponse>;
  subscribeStrategiesEvents: handleServerStreamingCall<SubscribeStrategiesEventsRequest, StrategyEvent>;
}

export interface TradeClient extends Client {
  start(
    request: StartRequest,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  start(
    request: StartRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  start(
    request: StartRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  stop(request: StopRequest, callback: (error: ServiceError | null, response: StopResponse) => void): ClientUnaryCall;
  stop(
    request: StopRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: StopResponse) => void,
  ): ClientUnaryCall;
  stop(
    request: StopRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: StopResponse) => void,
  ): ClientUnaryCall;
  isStarted(
    request: StartRequest,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  isStarted(
    request: StartRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  isStarted(
    request: StartRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: StartResponse) => void,
  ): ClientUnaryCall;
  changeConfig(
    request: ChangeConfigRequest,
    callback: (error: ServiceError | null, response: ChangeConfigResponse) => void,
  ): ClientUnaryCall;
  changeConfig(
    request: ChangeConfigRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: ChangeConfigResponse) => void,
  ): ClientUnaryCall;
  changeConfig(
    request: ChangeConfigRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: ChangeConfigResponse) => void,
  ): ClientUnaryCall;
  getConfig(
    request: GetConfigRequest,
    callback: (error: ServiceError | null, response: GetConfigResponse) => void,
  ): ClientUnaryCall;
  getConfig(
    request: GetConfigRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: GetConfigResponse) => void,
  ): ClientUnaryCall;
  getConfig(
    request: GetConfigRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: GetConfigResponse) => void,
  ): ClientUnaryCall;
  subscribeStrategiesEvents(
    request: SubscribeStrategiesEventsRequest,
    options?: Partial<CallOptions>,
  ): ClientReadableStream<StrategyEvent>;
  subscribeStrategiesEvents(
    request: SubscribeStrategiesEventsRequest,
    metadata?: Metadata,
    options?: Partial<CallOptions>,
  ): ClientReadableStream<StrategyEvent>;
}

export const TradeClient = makeGenericClientConstructor(TradeService, "trade.Trade") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): TradeClient;
  service: typeof TradeService;
  serviceName: string;
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
