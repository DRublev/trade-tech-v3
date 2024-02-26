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
import Long from "long";
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

export interface Quant {
  units: number;
  nano: number;
}

export interface OHLC {
  open: Quant | undefined;
  high: Quant | undefined;
  low: Quant | undefined;
  close: Quant | undefined;
  volume: number;
  time: Date | undefined;
}

/** Что хотим вернуть */
export interface GetCandlesResponse {
  candles: OHLC[];
}

export interface SubscribeCandlesRequest {
  instrumentId: string;
  interval: number;
}

export interface BidAsk {
  price: Quant | undefined;
  quantity: number;
}

export interface Orderbook {
  instrumentId: string;
  depth: number;
  time: Date | undefined;
  limitUp: Quant | undefined;
  limitDown: Quant | undefined;
  bids: BidAsk[];
  asks: BidAsk[];
}

export interface SubscribeOrderbookRequest {
  instrumentId: string;
  depth: number;
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

function createBaseQuant(): Quant {
  return { units: 0, nano: 0 };
}

export const Quant = {
  encode(message: Quant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.units !== 0) {
      writer.uint32(8).int32(message.units);
    }
    if (message.nano !== 0) {
      writer.uint32(16).int32(message.nano);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Quant {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuant();
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

  fromJSON(object: any): Quant {
    return {
      units: isSet(object.units) ? globalThis.Number(object.units) : 0,
      nano: isSet(object.nano) ? globalThis.Number(object.nano) : 0,
    };
  },

  toJSON(message: Quant): unknown {
    const obj: any = {};
    if (message.units !== 0) {
      obj.units = Math.round(message.units);
    }
    if (message.nano !== 0) {
      obj.nano = Math.round(message.nano);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Quant>, I>>(base?: I): Quant {
    return Quant.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Quant>, I>>(object: I): Quant {
    const message = createBaseQuant();
    message.units = object.units ?? 0;
    message.nano = object.nano ?? 0;
    return message;
  },
};

function createBaseOHLC(): OHLC {
  return { open: undefined, high: undefined, low: undefined, close: undefined, volume: 0, time: undefined };
}

export const OHLC = {
  encode(message: OHLC, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.open !== undefined) {
      Quant.encode(message.open, writer.uint32(10).fork()).ldelim();
    }
    if (message.high !== undefined) {
      Quant.encode(message.high, writer.uint32(18).fork()).ldelim();
    }
    if (message.low !== undefined) {
      Quant.encode(message.low, writer.uint32(26).fork()).ldelim();
    }
    if (message.close !== undefined) {
      Quant.encode(message.close, writer.uint32(34).fork()).ldelim();
    }
    if (message.volume !== 0) {
      writer.uint32(40).int64(message.volume);
    }
    if (message.time !== undefined) {
      Timestamp.encode(toTimestamp(message.time), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OHLC {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOHLC();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.open = Quant.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.high = Quant.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.low = Quant.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.close = Quant.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.volume = longToNumber(reader.int64() as Long);
          continue;
        case 6:
          if (tag !== 50) {
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

  fromJSON(object: any): OHLC {
    return {
      open: isSet(object.open) ? Quant.fromJSON(object.open) : undefined,
      high: isSet(object.high) ? Quant.fromJSON(object.high) : undefined,
      low: isSet(object.low) ? Quant.fromJSON(object.low) : undefined,
      close: isSet(object.close) ? Quant.fromJSON(object.close) : undefined,
      volume: isSet(object.volume) ? globalThis.Number(object.volume) : 0,
      time: isSet(object.time) ? fromJsonTimestamp(object.time) : undefined,
    };
  },

  toJSON(message: OHLC): unknown {
    const obj: any = {};
    if (message.open !== undefined) {
      obj.open = Quant.toJSON(message.open);
    }
    if (message.high !== undefined) {
      obj.high = Quant.toJSON(message.high);
    }
    if (message.low !== undefined) {
      obj.low = Quant.toJSON(message.low);
    }
    if (message.close !== undefined) {
      obj.close = Quant.toJSON(message.close);
    }
    if (message.volume !== 0) {
      obj.volume = Math.round(message.volume);
    }
    if (message.time !== undefined) {
      obj.time = message.time.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<OHLC>, I>>(base?: I): OHLC {
    return OHLC.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<OHLC>, I>>(object: I): OHLC {
    const message = createBaseOHLC();
    message.open = (object.open !== undefined && object.open !== null) ? Quant.fromPartial(object.open) : undefined;
    message.high = (object.high !== undefined && object.high !== null) ? Quant.fromPartial(object.high) : undefined;
    message.low = (object.low !== undefined && object.low !== null) ? Quant.fromPartial(object.low) : undefined;
    message.close = (object.close !== undefined && object.close !== null) ? Quant.fromPartial(object.close) : undefined;
    message.volume = object.volume ?? 0;
    message.time = object.time ?? undefined;
    return message;
  },
};

function createBaseGetCandlesResponse(): GetCandlesResponse {
  return { candles: [] };
}

export const GetCandlesResponse = {
  encode(message: GetCandlesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.candles) {
      OHLC.encode(v!, writer.uint32(10).fork()).ldelim();
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

          message.candles.push(OHLC.decode(reader, reader.uint32()));
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
      candles: globalThis.Array.isArray(object?.candles) ? object.candles.map((e: any) => OHLC.fromJSON(e)) : [],
    };
  },

  toJSON(message: GetCandlesResponse): unknown {
    const obj: any = {};
    if (message.candles?.length) {
      obj.candles = message.candles.map((e) => OHLC.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCandlesResponse>, I>>(base?: I): GetCandlesResponse {
    return GetCandlesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCandlesResponse>, I>>(object: I): GetCandlesResponse {
    const message = createBaseGetCandlesResponse();
    message.candles = object.candles?.map((e) => OHLC.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSubscribeCandlesRequest(): SubscribeCandlesRequest {
  return { instrumentId: "", interval: 0 };
}

export const SubscribeCandlesRequest = {
  encode(message: SubscribeCandlesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instrumentId !== "") {
      writer.uint32(10).string(message.instrumentId);
    }
    if (message.interval !== 0) {
      writer.uint32(16).int32(message.interval);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeCandlesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubscribeCandlesRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubscribeCandlesRequest {
    return {
      instrumentId: isSet(object.instrumentId) ? globalThis.String(object.instrumentId) : "",
      interval: isSet(object.interval) ? globalThis.Number(object.interval) : 0,
    };
  },

  toJSON(message: SubscribeCandlesRequest): unknown {
    const obj: any = {};
    if (message.instrumentId !== "") {
      obj.instrumentId = message.instrumentId;
    }
    if (message.interval !== 0) {
      obj.interval = Math.round(message.interval);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubscribeCandlesRequest>, I>>(base?: I): SubscribeCandlesRequest {
    return SubscribeCandlesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubscribeCandlesRequest>, I>>(object: I): SubscribeCandlesRequest {
    const message = createBaseSubscribeCandlesRequest();
    message.instrumentId = object.instrumentId ?? "";
    message.interval = object.interval ?? 0;
    return message;
  },
};

function createBaseBidAsk(): BidAsk {
  return { price: undefined, quantity: 0 };
}

export const BidAsk = {
  encode(message: BidAsk, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.price !== undefined) {
      Quant.encode(message.price, writer.uint32(10).fork()).ldelim();
    }
    if (message.quantity !== 0) {
      writer.uint32(16).int64(message.quantity);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BidAsk {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBidAsk();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.price = Quant.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.quantity = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): BidAsk {
    return {
      price: isSet(object.price) ? Quant.fromJSON(object.price) : undefined,
      quantity: isSet(object.quantity) ? globalThis.Number(object.quantity) : 0,
    };
  },

  toJSON(message: BidAsk): unknown {
    const obj: any = {};
    if (message.price !== undefined) {
      obj.price = Quant.toJSON(message.price);
    }
    if (message.quantity !== 0) {
      obj.quantity = Math.round(message.quantity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BidAsk>, I>>(base?: I): BidAsk {
    return BidAsk.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BidAsk>, I>>(object: I): BidAsk {
    const message = createBaseBidAsk();
    message.price = (object.price !== undefined && object.price !== null) ? Quant.fromPartial(object.price) : undefined;
    message.quantity = object.quantity ?? 0;
    return message;
  },
};

function createBaseOrderbook(): Orderbook {
  return { instrumentId: "", depth: 0, time: undefined, limitUp: undefined, limitDown: undefined, bids: [], asks: [] };
}

export const Orderbook = {
  encode(message: Orderbook, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instrumentId !== "") {
      writer.uint32(10).string(message.instrumentId);
    }
    if (message.depth !== 0) {
      writer.uint32(16).int32(message.depth);
    }
    if (message.time !== undefined) {
      Timestamp.encode(toTimestamp(message.time), writer.uint32(26).fork()).ldelim();
    }
    if (message.limitUp !== undefined) {
      Quant.encode(message.limitUp, writer.uint32(34).fork()).ldelim();
    }
    if (message.limitDown !== undefined) {
      Quant.encode(message.limitDown, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.bids) {
      BidAsk.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.asks) {
      BidAsk.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Orderbook {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderbook();
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

          message.depth = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.time = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.limitUp = Quant.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.limitDown = Quant.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.bids.push(BidAsk.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.asks.push(BidAsk.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Orderbook {
    return {
      instrumentId: isSet(object.instrumentId) ? globalThis.String(object.instrumentId) : "",
      depth: isSet(object.depth) ? globalThis.Number(object.depth) : 0,
      time: isSet(object.time) ? fromJsonTimestamp(object.time) : undefined,
      limitUp: isSet(object.limitUp) ? Quant.fromJSON(object.limitUp) : undefined,
      limitDown: isSet(object.limitDown) ? Quant.fromJSON(object.limitDown) : undefined,
      bids: globalThis.Array.isArray(object?.bids) ? object.bids.map((e: any) => BidAsk.fromJSON(e)) : [],
      asks: globalThis.Array.isArray(object?.asks) ? object.asks.map((e: any) => BidAsk.fromJSON(e)) : [],
    };
  },

  toJSON(message: Orderbook): unknown {
    const obj: any = {};
    if (message.instrumentId !== "") {
      obj.instrumentId = message.instrumentId;
    }
    if (message.depth !== 0) {
      obj.depth = Math.round(message.depth);
    }
    if (message.time !== undefined) {
      obj.time = message.time.toISOString();
    }
    if (message.limitUp !== undefined) {
      obj.limitUp = Quant.toJSON(message.limitUp);
    }
    if (message.limitDown !== undefined) {
      obj.limitDown = Quant.toJSON(message.limitDown);
    }
    if (message.bids?.length) {
      obj.bids = message.bids.map((e) => BidAsk.toJSON(e));
    }
    if (message.asks?.length) {
      obj.asks = message.asks.map((e) => BidAsk.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Orderbook>, I>>(base?: I): Orderbook {
    return Orderbook.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Orderbook>, I>>(object: I): Orderbook {
    const message = createBaseOrderbook();
    message.instrumentId = object.instrumentId ?? "";
    message.depth = object.depth ?? 0;
    message.time = object.time ?? undefined;
    message.limitUp = (object.limitUp !== undefined && object.limitUp !== null)
      ? Quant.fromPartial(object.limitUp)
      : undefined;
    message.limitDown = (object.limitDown !== undefined && object.limitDown !== null)
      ? Quant.fromPartial(object.limitDown)
      : undefined;
    message.bids = object.bids?.map((e) => BidAsk.fromPartial(e)) || [];
    message.asks = object.asks?.map((e) => BidAsk.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSubscribeOrderbookRequest(): SubscribeOrderbookRequest {
  return { instrumentId: "", depth: 0 };
}

export const SubscribeOrderbookRequest = {
  encode(message: SubscribeOrderbookRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instrumentId !== "") {
      writer.uint32(10).string(message.instrumentId);
    }
    if (message.depth !== 0) {
      writer.uint32(16).int32(message.depth);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubscribeOrderbookRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubscribeOrderbookRequest();
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

          message.depth = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SubscribeOrderbookRequest {
    return {
      instrumentId: isSet(object.instrumentId) ? globalThis.String(object.instrumentId) : "",
      depth: isSet(object.depth) ? globalThis.Number(object.depth) : 0,
    };
  },

  toJSON(message: SubscribeOrderbookRequest): unknown {
    const obj: any = {};
    if (message.instrumentId !== "") {
      obj.instrumentId = message.instrumentId;
    }
    if (message.depth !== 0) {
      obj.depth = Math.round(message.depth);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SubscribeOrderbookRequest>, I>>(base?: I): SubscribeOrderbookRequest {
    return SubscribeOrderbookRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SubscribeOrderbookRequest>, I>>(object: I): SubscribeOrderbookRequest {
    const message = createBaseSubscribeOrderbookRequest();
    message.instrumentId = object.instrumentId ?? "";
    message.depth = object.depth ?? 0;
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
  subscribeCandles: {
    path: "/marketData.MarketData/SubscribeCandles",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: SubscribeCandlesRequest) => Buffer.from(SubscribeCandlesRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SubscribeCandlesRequest.decode(value),
    responseSerialize: (value: OHLC) => Buffer.from(OHLC.encode(value).finish()),
    responseDeserialize: (value: Buffer) => OHLC.decode(value),
  },
  subscribeOrderbook: {
    path: "/marketData.MarketData/SubscribeOrderbook",
    requestStream: false,
    responseStream: true,
    requestSerialize: (value: SubscribeOrderbookRequest) =>
      Buffer.from(SubscribeOrderbookRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SubscribeOrderbookRequest.decode(value),
    responseSerialize: (value: Orderbook) => Buffer.from(Orderbook.encode(value).finish()),
    responseDeserialize: (value: Buffer) => Orderbook.decode(value),
  },
} as const;

export interface MarketDataServer extends UntypedServiceImplementation {
  /** Название нашего эндпоинта */
  getCandles: handleUnaryCall<GetCandlesRequest, GetCandlesResponse>;
  subscribeCandles: handleServerStreamingCall<SubscribeCandlesRequest, OHLC>;
  subscribeOrderbook: handleServerStreamingCall<SubscribeOrderbookRequest, Orderbook>;
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
  subscribeCandles(request: SubscribeCandlesRequest, options?: Partial<CallOptions>): ClientReadableStream<OHLC>;
  subscribeCandles(
    request: SubscribeCandlesRequest,
    metadata?: Metadata,
    options?: Partial<CallOptions>,
  ): ClientReadableStream<OHLC>;
  subscribeOrderbook(
    request: SubscribeOrderbookRequest,
    options?: Partial<CallOptions>,
  ): ClientReadableStream<Orderbook>;
  subscribeOrderbook(
    request: SubscribeOrderbookRequest,
    metadata?: Metadata,
    options?: Partial<CallOptions>,
  ): ClientReadableStream<Orderbook>;
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

function longToNumber(long: Long): number {
  if (long.gt(globalThis.Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
