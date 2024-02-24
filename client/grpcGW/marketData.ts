/* eslint-disable */
import { ChannelCredentials, Client, makeGenericClientConstructor, Metadata } from "@grpc/grpc-js";
import type {
  CallOptions,
  ClientOptions,
  ClientUnaryCall,
  handleUnaryCall,
  ServiceError,
  UntypedServiceImplementation,
} from "@grpc/grpc-js";
import _m0 from "protobufjs/minimal";
import { Timestamp } from "./google/protobuf/timestamp";

export const protobufPackage = "marketData";

/** Что ждем от клиента */
export interface GetCandlesRequest {
  instrumentId: string;
  interval: number;
  start: Date | undefined;
  end: Date | undefined;
}

/** Что хотим вернуть */
export interface GetCandlesResponse {
  candles: GetCandlesResponse_OHLC[];
}

export interface GetCandlesResponse_Quant {
  units: number;
  nano: number;
}

export interface GetCandlesResponse_OHLC {
  open: GetCandlesResponse_Quant | undefined;
  high: GetCandlesResponse_Quant | undefined;
  low: GetCandlesResponse_Quant | undefined;
  close: GetCandlesResponse_Quant | undefined;
  time: Date | undefined;
}

function createBaseGetCandlesRequest(): GetCandlesRequest {
  return { instrumentId: "", interval: 0, start: undefined, end: undefined };
}

export const GetCandlesRequest = {
  encode(message: GetCandlesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instrumentId !== "") {
      writer.uint32(10).string(message.instrumentId);
    }
    if (message.interval !== 0) {
      writer.uint32(16).int32(message.interval);
    }
    if (message.start !== undefined) {
      Timestamp.encode(toTimestamp(message.start), writer.uint32(26).fork()).ldelim();
    }
    if (message.end !== undefined) {
      Timestamp.encode(toTimestamp(message.end), writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCandlesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCandlesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.instrumentId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.interval = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.start = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.end = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCandlesRequest {
    return {
      instrumentId: isSet(object.instrumentId) ? globalThis.String(object.instrumentId) : "",
      interval: isSet(object.interval) ? globalThis.Number(object.interval) : 0,
      start: isSet(object.start) ? fromJsonTimestamp(object.start) : undefined,
      end: isSet(object.end) ? fromJsonTimestamp(object.end) : undefined,
    };
  },

  toJSON(message: GetCandlesRequest): unknown {
    const obj: any = {};
    if (message.instrumentId !== "") {
      obj.instrumentId = message.instrumentId;
    }
    if (message.interval !== 0) {
      obj.interval = Math.round(message.interval);
    }
    if (message.start !== undefined) {
      obj.start = message.start.toISOString();
    }
    if (message.end !== undefined) {
      obj.end = message.end.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCandlesRequest>, I>>(base?: I): GetCandlesRequest {
    return GetCandlesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCandlesRequest>, I>>(object: I): GetCandlesRequest {
    const message = createBaseGetCandlesRequest();
    message.instrumentId = object.instrumentId ?? "";
    message.interval = object.interval ?? 0;
    message.start = object.start ?? undefined;
    message.end = object.end ?? undefined;
    return message;
  },
};

function createBaseGetCandlesResponse(): GetCandlesResponse {
  return { candles: [] };
}

export const GetCandlesResponse = {
  encode(message: GetCandlesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.candles) {
      GetCandlesResponse_OHLC.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCandlesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCandlesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.candles.push(GetCandlesResponse_OHLC.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCandlesResponse {
    return {
      candles: globalThis.Array.isArray(object?.candles)
        ? object.candles.map((e: any) => GetCandlesResponse_OHLC.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetCandlesResponse): unknown {
    const obj: any = {};
    if (message.candles?.length) {
      obj.candles = message.candles.map((e) => GetCandlesResponse_OHLC.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCandlesResponse>, I>>(base?: I): GetCandlesResponse {
    return GetCandlesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCandlesResponse>, I>>(object: I): GetCandlesResponse {
    const message = createBaseGetCandlesResponse();
    message.candles = object.candles?.map((e) => GetCandlesResponse_OHLC.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetCandlesResponse_Quant(): GetCandlesResponse_Quant {
  return { units: 0, nano: 0 };
}

export const GetCandlesResponse_Quant = {
  encode(message: GetCandlesResponse_Quant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.units !== 0) {
      writer.uint32(8).int32(message.units);
    }
    if (message.nano !== 0) {
      writer.uint32(16).int32(message.nano);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCandlesResponse_Quant {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCandlesResponse_Quant();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.units = reader.int32();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.nano = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCandlesResponse_Quant {
    return {
      units: isSet(object.units) ? globalThis.Number(object.units) : 0,
      nano: isSet(object.nano) ? globalThis.Number(object.nano) : 0,
    };
  },

  toJSON(message: GetCandlesResponse_Quant): unknown {
    const obj: any = {};
    if (message.units !== 0) {
      obj.units = Math.round(message.units);
    }
    if (message.nano !== 0) {
      obj.nano = Math.round(message.nano);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCandlesResponse_Quant>, I>>(base?: I): GetCandlesResponse_Quant {
    return GetCandlesResponse_Quant.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCandlesResponse_Quant>, I>>(object: I): GetCandlesResponse_Quant {
    const message = createBaseGetCandlesResponse_Quant();
    message.units = object.units ?? 0;
    message.nano = object.nano ?? 0;
    return message;
  },
};

function createBaseGetCandlesResponse_OHLC(): GetCandlesResponse_OHLC {
  return { open: undefined, high: undefined, low: undefined, close: undefined, time: undefined };
}

export const GetCandlesResponse_OHLC = {
  encode(message: GetCandlesResponse_OHLC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.open !== undefined) {
      GetCandlesResponse_Quant.encode(message.open, writer.uint32(10).fork()).ldelim();
    }
    if (message.high !== undefined) {
      GetCandlesResponse_Quant.encode(message.high, writer.uint32(18).fork()).ldelim();
    }
    if (message.low !== undefined) {
      GetCandlesResponse_Quant.encode(message.low, writer.uint32(26).fork()).ldelim();
    }
    if (message.close !== undefined) {
      GetCandlesResponse_Quant.encode(message.close, writer.uint32(34).fork()).ldelim();
    }
    if (message.time !== undefined) {
      Timestamp.encode(toTimestamp(message.time), writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCandlesResponse_OHLC {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCandlesResponse_OHLC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.open = GetCandlesResponse_Quant.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.high = GetCandlesResponse_Quant.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.low = GetCandlesResponse_Quant.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.close = GetCandlesResponse_Quant.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.time = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCandlesResponse_OHLC {
    return {
      open: isSet(object.open) ? GetCandlesResponse_Quant.fromJSON(object.open) : undefined,
      high: isSet(object.high) ? GetCandlesResponse_Quant.fromJSON(object.high) : undefined,
      low: isSet(object.low) ? GetCandlesResponse_Quant.fromJSON(object.low) : undefined,
      close: isSet(object.close) ? GetCandlesResponse_Quant.fromJSON(object.close) : undefined,
      time: isSet(object.time) ? fromJsonTimestamp(object.time) : undefined,
    };
  },

  toJSON(message: GetCandlesResponse_OHLC): unknown {
    const obj: any = {};
    if (message.open !== undefined) {
      obj.open = GetCandlesResponse_Quant.toJSON(message.open);
    }
    if (message.high !== undefined) {
      obj.high = GetCandlesResponse_Quant.toJSON(message.high);
    }
    if (message.low !== undefined) {
      obj.low = GetCandlesResponse_Quant.toJSON(message.low);
    }
    if (message.close !== undefined) {
      obj.close = GetCandlesResponse_Quant.toJSON(message.close);
    }
    if (message.time !== undefined) {
      obj.time = message.time.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCandlesResponse_OHLC>, I>>(base?: I): GetCandlesResponse_OHLC {
    return GetCandlesResponse_OHLC.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCandlesResponse_OHLC>, I>>(object: I): GetCandlesResponse_OHLC {
    const message = createBaseGetCandlesResponse_OHLC();
    message.open = (object.open !== undefined && object.open !== null)
      ? GetCandlesResponse_Quant.fromPartial(object.open)
      : undefined;
    message.high = (object.high !== undefined && object.high !== null)
      ? GetCandlesResponse_Quant.fromPartial(object.high)
      : undefined;
    message.low = (object.low !== undefined && object.low !== null)
      ? GetCandlesResponse_Quant.fromPartial(object.low)
      : undefined;
    message.close = (object.close !== undefined && object.close !== null)
      ? GetCandlesResponse_Quant.fromPartial(object.close)
      : undefined;
    message.time = object.time ?? undefined;
    return message;
  },
};

export type MarketDataService = typeof MarketDataService;
export const MarketDataService = {
  /** Название нашего эндпоинта */
  getCandles: {
    path: "/marketData.MarketData/GetCandles",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: GetCandlesRequest) => Buffer.from(GetCandlesRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => GetCandlesRequest.decode(value),
    responseSerialize: (value: GetCandlesResponse) => Buffer.from(GetCandlesResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => GetCandlesResponse.decode(value),
  },
} as const;

export interface MarketDataServer extends UntypedServiceImplementation {
  /** Название нашего эндпоинта */
  getCandles: handleUnaryCall<GetCandlesRequest, GetCandlesResponse>;
}

export interface MarketDataClient extends Client {
  /** Название нашего эндпоинта */
  getCandles(
    request: GetCandlesRequest,
    callback: (error: ServiceError | null, response: GetCandlesResponse) => void,
  ): ClientUnaryCall;
  getCandles(
    request: GetCandlesRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: GetCandlesResponse) => void,
  ): ClientUnaryCall;
  getCandles(
    request: GetCandlesRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: GetCandlesResponse) => void,
  ): ClientUnaryCall;
}

export const MarketDataClient = makeGenericClientConstructor(MarketDataService, "marketData.MarketData") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): MarketDataClient;
  service: typeof MarketDataService;
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

function toTimestamp(date: Date): Timestamp {
  const seconds = Math.trunc(date.getTime() / 1_000);
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = (t.seconds || 0) * 1_000;
  millis += (t.nanos || 0) / 1_000_000;
  return new globalThis.Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof globalThis.Date) {
    return o;
  } else if (typeof o === "string") {
    return new globalThis.Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
