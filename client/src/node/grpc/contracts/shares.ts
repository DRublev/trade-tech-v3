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

export const protobufPackage = "shares";

export interface Quatation {
  units: number;
  nano: number;
}

export interface Share {
  figi: string;
  name: string;
  exchange: string;
  ticker: string;
  lot: number;
  ipoDate: Date | undefined;
  tradingStatus: number;
  minPriceIncrement: Quatation | undefined;
  uid: string;
  first1minCandleDate: Date | undefined;
  first1dayCandleDate: Date | undefined;
}

export interface GetInstrumentsRequest {
  instrumentStatus: number;
}

export interface GetSharesResponse {
  instruments: Share[];
}

export interface GetTradingSchedulesRequest {
  exchange: string;
  from: Date | undefined;
  to: Date | undefined;
}

export interface GetTradingSchedulesResponse {
  exchanges: TradingSchedule[];
}

export interface TradingSchedule {
  exchange: string;
  days: TradingDay[];
}

export interface TradingDay {
  date: Date | undefined;
  isTradingDay: boolean;
  startTime: Date | undefined;
  endTime: Date | undefined;
  openingAuctionStartTime: Date | undefined;
  closingAuctionEndTime: Date | undefined;
  eveningOpeningAuctionStartTime: Date | undefined;
  eveningStartTime: Date | undefined;
  eveningEndTime: Date | undefined;
  clearingStartTime: Date | undefined;
  clearingEndTime: Date | undefined;
  premarketStartTime: Date | undefined;
  premarketEndTime: Date | undefined;
  closingAuctionStartTime: Date | undefined;
  peningAuctionEndTime: Date | undefined;
}

function createBaseQuatation(): Quatation {
  return { units: 0, nano: 0 };
}

export const Quatation = {
  encode(message: Quatation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.units !== 0) {
      writer.uint32(8).int32(message.units);
    }
    if (message.nano !== 0) {
      writer.uint32(16).int32(message.nano);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Quatation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuatation();
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

  fromJSON(object: any): Quatation {
    return {
      units: isSet(object.units) ? globalThis.Number(object.units) : 0,
      nano: isSet(object.nano) ? globalThis.Number(object.nano) : 0,
    };
  },

  toJSON(message: Quatation): unknown {
    const obj: any = {};
    if (message.units !== 0) {
      obj.units = Math.round(message.units);
    }
    if (message.nano !== 0) {
      obj.nano = Math.round(message.nano);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Quatation>, I>>(base?: I): Quatation {
    return Quatation.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Quatation>, I>>(object: I): Quatation {
    const message = createBaseQuatation();
    message.units = object.units ?? 0;
    message.nano = object.nano ?? 0;
    return message;
  },
};

function createBaseShare(): Share {
  return {
    figi: "",
    name: "",
    exchange: "",
    ticker: "",
    lot: 0,
    ipoDate: undefined,
    tradingStatus: 0,
    minPriceIncrement: undefined,
    uid: "",
    first1minCandleDate: undefined,
    first1dayCandleDate: undefined,
  };
}

export const Share = {
  encode(message: Share, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.figi !== "") {
      writer.uint32(10).string(message.figi);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.exchange !== "") {
      writer.uint32(26).string(message.exchange);
    }
    if (message.ticker !== "") {
      writer.uint32(34).string(message.ticker);
    }
    if (message.lot !== 0) {
      writer.uint32(40).int32(message.lot);
    }
    if (message.ipoDate !== undefined) {
      Timestamp.encode(toTimestamp(message.ipoDate), writer.uint32(50).fork()).ldelim();
    }
    if (message.tradingStatus !== 0) {
      writer.uint32(56).int32(message.tradingStatus);
    }
    if (message.minPriceIncrement !== undefined) {
      Quatation.encode(message.minPriceIncrement, writer.uint32(66).fork()).ldelim();
    }
    if (message.uid !== "") {
      writer.uint32(74).string(message.uid);
    }
    if (message.first1minCandleDate !== undefined) {
      Timestamp.encode(toTimestamp(message.first1minCandleDate), writer.uint32(82).fork()).ldelim();
    }
    if (message.first1dayCandleDate !== undefined) {
      Timestamp.encode(toTimestamp(message.first1dayCandleDate), writer.uint32(90).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Share {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseShare();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.figi = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.exchange = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.ticker = reader.string();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.lot = reader.int32();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.ipoDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.tradingStatus = reader.int32();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.minPriceIncrement = Quatation.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.uid = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.first1minCandleDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.first1dayCandleDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Share {
    return {
      figi: isSet(object.figi) ? globalThis.String(object.figi) : "",
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      exchange: isSet(object.exchange) ? globalThis.String(object.exchange) : "",
      ticker: isSet(object.ticker) ? globalThis.String(object.ticker) : "",
      lot: isSet(object.lot) ? globalThis.Number(object.lot) : 0,
      ipoDate: isSet(object.ipoDate) ? fromJsonTimestamp(object.ipoDate) : undefined,
      tradingStatus: isSet(object.tradingStatus) ? globalThis.Number(object.tradingStatus) : 0,
      minPriceIncrement: isSet(object.minPriceIncrement) ? Quatation.fromJSON(object.minPriceIncrement) : undefined,
      uid: isSet(object.uid) ? globalThis.String(object.uid) : "",
      first1minCandleDate: isSet(object.first1minCandleDate)
        ? fromJsonTimestamp(object.first1minCandleDate)
        : undefined,
      first1dayCandleDate: isSet(object.first1dayCandleDate)
        ? fromJsonTimestamp(object.first1dayCandleDate)
        : undefined,
    };
  },

  toJSON(message: Share): unknown {
    const obj: any = {};
    if (message.figi !== "") {
      obj.figi = message.figi;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.exchange !== "") {
      obj.exchange = message.exchange;
    }
    if (message.ticker !== "") {
      obj.ticker = message.ticker;
    }
    if (message.lot !== 0) {
      obj.lot = Math.round(message.lot);
    }
    if (message.ipoDate !== undefined) {
      obj.ipoDate = message.ipoDate.toISOString();
    }
    if (message.tradingStatus !== 0) {
      obj.tradingStatus = Math.round(message.tradingStatus);
    }
    if (message.minPriceIncrement !== undefined) {
      obj.minPriceIncrement = Quatation.toJSON(message.minPriceIncrement);
    }
    if (message.uid !== "") {
      obj.uid = message.uid;
    }
    if (message.first1minCandleDate !== undefined) {
      obj.first1minCandleDate = message.first1minCandleDate.toISOString();
    }
    if (message.first1dayCandleDate !== undefined) {
      obj.first1dayCandleDate = message.first1dayCandleDate.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Share>, I>>(base?: I): Share {
    return Share.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Share>, I>>(object: I): Share {
    const message = createBaseShare();
    message.figi = object.figi ?? "";
    message.name = object.name ?? "";
    message.exchange = object.exchange ?? "";
    message.ticker = object.ticker ?? "";
    message.lot = object.lot ?? 0;
    message.ipoDate = object.ipoDate ?? undefined;
    message.tradingStatus = object.tradingStatus ?? 0;
    message.minPriceIncrement = (object.minPriceIncrement !== undefined && object.minPriceIncrement !== null)
      ? Quatation.fromPartial(object.minPriceIncrement)
      : undefined;
    message.uid = object.uid ?? "";
    message.first1minCandleDate = object.first1minCandleDate ?? undefined;
    message.first1dayCandleDate = object.first1dayCandleDate ?? undefined;
    return message;
  },
};

function createBaseGetInstrumentsRequest(): GetInstrumentsRequest {
  return { instrumentStatus: 0 };
}

export const GetInstrumentsRequest = {
  encode(message: GetInstrumentsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.instrumentStatus !== 0) {
      writer.uint32(16).int32(message.instrumentStatus);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetInstrumentsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetInstrumentsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag !== 16) {
            break;
          }

          message.instrumentStatus = reader.int32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetInstrumentsRequest {
    return { instrumentStatus: isSet(object.instrumentStatus) ? globalThis.Number(object.instrumentStatus) : 0 };
  },

  toJSON(message: GetInstrumentsRequest): unknown {
    const obj: any = {};
    if (message.instrumentStatus !== 0) {
      obj.instrumentStatus = Math.round(message.instrumentStatus);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetInstrumentsRequest>, I>>(base?: I): GetInstrumentsRequest {
    return GetInstrumentsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetInstrumentsRequest>, I>>(object: I): GetInstrumentsRequest {
    const message = createBaseGetInstrumentsRequest();
    message.instrumentStatus = object.instrumentStatus ?? 0;
    return message;
  },
};

function createBaseGetSharesResponse(): GetSharesResponse {
  return { instruments: [] };
}

export const GetSharesResponse = {
  encode(message: GetSharesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.instruments) {
      Share.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetSharesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetSharesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.instruments.push(Share.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetSharesResponse {
    return {
      instruments: globalThis.Array.isArray(object?.instruments)
        ? object.instruments.map((e: any) => Share.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetSharesResponse): unknown {
    const obj: any = {};
    if (message.instruments?.length) {
      obj.instruments = message.instruments.map((e) => Share.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetSharesResponse>, I>>(base?: I): GetSharesResponse {
    return GetSharesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetSharesResponse>, I>>(object: I): GetSharesResponse {
    const message = createBaseGetSharesResponse();
    message.instruments = object.instruments?.map((e) => Share.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetTradingSchedulesRequest(): GetTradingSchedulesRequest {
  return { exchange: "", from: undefined, to: undefined };
}

export const GetTradingSchedulesRequest = {
  encode(message: GetTradingSchedulesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exchange !== "") {
      writer.uint32(10).string(message.exchange);
    }
    if (message.from !== undefined) {
      Timestamp.encode(toTimestamp(message.from), writer.uint32(26).fork()).ldelim();
    }
    if (message.to !== undefined) {
      Timestamp.encode(toTimestamp(message.to), writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTradingSchedulesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTradingSchedulesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.exchange = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.from = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.to = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetTradingSchedulesRequest {
    return {
      exchange: isSet(object.exchange) ? globalThis.String(object.exchange) : "",
      from: isSet(object.from) ? fromJsonTimestamp(object.from) : undefined,
      to: isSet(object.to) ? fromJsonTimestamp(object.to) : undefined,
    };
  },

  toJSON(message: GetTradingSchedulesRequest): unknown {
    const obj: any = {};
    if (message.exchange !== "") {
      obj.exchange = message.exchange;
    }
    if (message.from !== undefined) {
      obj.from = message.from.toISOString();
    }
    if (message.to !== undefined) {
      obj.to = message.to.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetTradingSchedulesRequest>, I>>(base?: I): GetTradingSchedulesRequest {
    return GetTradingSchedulesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetTradingSchedulesRequest>, I>>(object: I): GetTradingSchedulesRequest {
    const message = createBaseGetTradingSchedulesRequest();
    message.exchange = object.exchange ?? "";
    message.from = object.from ?? undefined;
    message.to = object.to ?? undefined;
    return message;
  },
};

function createBaseGetTradingSchedulesResponse(): GetTradingSchedulesResponse {
  return { exchanges: [] };
}

export const GetTradingSchedulesResponse = {
  encode(message: GetTradingSchedulesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.exchanges) {
      TradingSchedule.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTradingSchedulesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTradingSchedulesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.exchanges.push(TradingSchedule.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetTradingSchedulesResponse {
    return {
      exchanges: globalThis.Array.isArray(object?.exchanges)
        ? object.exchanges.map((e: any) => TradingSchedule.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetTradingSchedulesResponse): unknown {
    const obj: any = {};
    if (message.exchanges?.length) {
      obj.exchanges = message.exchanges.map((e) => TradingSchedule.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetTradingSchedulesResponse>, I>>(base?: I): GetTradingSchedulesResponse {
    return GetTradingSchedulesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetTradingSchedulesResponse>, I>>(object: I): GetTradingSchedulesResponse {
    const message = createBaseGetTradingSchedulesResponse();
    message.exchanges = object.exchanges?.map((e) => TradingSchedule.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTradingSchedule(): TradingSchedule {
  return { exchange: "", days: [] };
}

export const TradingSchedule = {
  encode(message: TradingSchedule, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exchange !== "") {
      writer.uint32(10).string(message.exchange);
    }
    for (const v of message.days) {
      TradingDay.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TradingSchedule {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTradingSchedule();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.exchange = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.days.push(TradingDay.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TradingSchedule {
    return {
      exchange: isSet(object.exchange) ? globalThis.String(object.exchange) : "",
      days: globalThis.Array.isArray(object?.days) ? object.days.map((e: any) => TradingDay.fromJSON(e)) : [],
    };
  },

  toJSON(message: TradingSchedule): unknown {
    const obj: any = {};
    if (message.exchange !== "") {
      obj.exchange = message.exchange;
    }
    if (message.days?.length) {
      obj.days = message.days.map((e) => TradingDay.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TradingSchedule>, I>>(base?: I): TradingSchedule {
    return TradingSchedule.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TradingSchedule>, I>>(object: I): TradingSchedule {
    const message = createBaseTradingSchedule();
    message.exchange = object.exchange ?? "";
    message.days = object.days?.map((e) => TradingDay.fromPartial(e)) || [];
    return message;
  },
};

function createBaseTradingDay(): TradingDay {
  return {
    date: undefined,
    isTradingDay: false,
    startTime: undefined,
    endTime: undefined,
    openingAuctionStartTime: undefined,
    closingAuctionEndTime: undefined,
    eveningOpeningAuctionStartTime: undefined,
    eveningStartTime: undefined,
    eveningEndTime: undefined,
    clearingStartTime: undefined,
    clearingEndTime: undefined,
    premarketStartTime: undefined,
    premarketEndTime: undefined,
    closingAuctionStartTime: undefined,
    peningAuctionEndTime: undefined,
  };
}

export const TradingDay = {
  encode(message: TradingDay, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.date !== undefined) {
      Timestamp.encode(toTimestamp(message.date), writer.uint32(10).fork()).ldelim();
    }
    if (message.isTradingDay !== false) {
      writer.uint32(16).bool(message.isTradingDay);
    }
    if (message.startTime !== undefined) {
      Timestamp.encode(toTimestamp(message.startTime), writer.uint32(26).fork()).ldelim();
    }
    if (message.endTime !== undefined) {
      Timestamp.encode(toTimestamp(message.endTime), writer.uint32(34).fork()).ldelim();
    }
    if (message.openingAuctionStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.openingAuctionStartTime), writer.uint32(42).fork()).ldelim();
    }
    if (message.closingAuctionEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.closingAuctionEndTime), writer.uint32(50).fork()).ldelim();
    }
    if (message.eveningOpeningAuctionStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.eveningOpeningAuctionStartTime), writer.uint32(58).fork()).ldelim();
    }
    if (message.eveningStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.eveningStartTime), writer.uint32(66).fork()).ldelim();
    }
    if (message.eveningEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.eveningEndTime), writer.uint32(74).fork()).ldelim();
    }
    if (message.clearingStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.clearingStartTime), writer.uint32(82).fork()).ldelim();
    }
    if (message.clearingEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.clearingEndTime), writer.uint32(90).fork()).ldelim();
    }
    if (message.premarketStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.premarketStartTime), writer.uint32(98).fork()).ldelim();
    }
    if (message.premarketEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.premarketEndTime), writer.uint32(106).fork()).ldelim();
    }
    if (message.closingAuctionStartTime !== undefined) {
      Timestamp.encode(toTimestamp(message.closingAuctionStartTime), writer.uint32(114).fork()).ldelim();
    }
    if (message.peningAuctionEndTime !== undefined) {
      Timestamp.encode(toTimestamp(message.peningAuctionEndTime), writer.uint32(122).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TradingDay {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTradingDay();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.date = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.isTradingDay = reader.bool();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.startTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.endTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.openingAuctionStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.closingAuctionEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.eveningOpeningAuctionStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.eveningStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.eveningEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.clearingStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.clearingEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.premarketStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.premarketEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.closingAuctionStartTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 15:
          if (tag !== 122) {
            break;
          }

          message.peningAuctionEndTime = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TradingDay {
    return {
      date: isSet(object.date) ? fromJsonTimestamp(object.date) : undefined,
      isTradingDay: isSet(object.isTradingDay) ? globalThis.Boolean(object.isTradingDay) : false,
      startTime: isSet(object.startTime) ? fromJsonTimestamp(object.startTime) : undefined,
      endTime: isSet(object.endTime) ? fromJsonTimestamp(object.endTime) : undefined,
      openingAuctionStartTime: isSet(object.openingAuctionStartTime)
        ? fromJsonTimestamp(object.openingAuctionStartTime)
        : undefined,
      closingAuctionEndTime: isSet(object.closingAuctionEndTime)
        ? fromJsonTimestamp(object.closingAuctionEndTime)
        : undefined,
      eveningOpeningAuctionStartTime: isSet(object.eveningOpeningAuctionStartTime)
        ? fromJsonTimestamp(object.eveningOpeningAuctionStartTime)
        : undefined,
      eveningStartTime: isSet(object.eveningStartTime) ? fromJsonTimestamp(object.eveningStartTime) : undefined,
      eveningEndTime: isSet(object.eveningEndTime) ? fromJsonTimestamp(object.eveningEndTime) : undefined,
      clearingStartTime: isSet(object.clearingStartTime) ? fromJsonTimestamp(object.clearingStartTime) : undefined,
      clearingEndTime: isSet(object.clearingEndTime) ? fromJsonTimestamp(object.clearingEndTime) : undefined,
      premarketStartTime: isSet(object.premarketStartTime) ? fromJsonTimestamp(object.premarketStartTime) : undefined,
      premarketEndTime: isSet(object.premarketEndTime) ? fromJsonTimestamp(object.premarketEndTime) : undefined,
      closingAuctionStartTime: isSet(object.closingAuctionStartTime)
        ? fromJsonTimestamp(object.closingAuctionStartTime)
        : undefined,
      peningAuctionEndTime: isSet(object.peningAuctionEndTime)
        ? fromJsonTimestamp(object.peningAuctionEndTime)
        : undefined,
    };
  },

  toJSON(message: TradingDay): unknown {
    const obj: any = {};
    if (message.date !== undefined) {
      obj.date = message.date.toISOString();
    }
    if (message.isTradingDay !== false) {
      obj.isTradingDay = message.isTradingDay;
    }
    if (message.startTime !== undefined) {
      obj.startTime = message.startTime.toISOString();
    }
    if (message.endTime !== undefined) {
      obj.endTime = message.endTime.toISOString();
    }
    if (message.openingAuctionStartTime !== undefined) {
      obj.openingAuctionStartTime = message.openingAuctionStartTime.toISOString();
    }
    if (message.closingAuctionEndTime !== undefined) {
      obj.closingAuctionEndTime = message.closingAuctionEndTime.toISOString();
    }
    if (message.eveningOpeningAuctionStartTime !== undefined) {
      obj.eveningOpeningAuctionStartTime = message.eveningOpeningAuctionStartTime.toISOString();
    }
    if (message.eveningStartTime !== undefined) {
      obj.eveningStartTime = message.eveningStartTime.toISOString();
    }
    if (message.eveningEndTime !== undefined) {
      obj.eveningEndTime = message.eveningEndTime.toISOString();
    }
    if (message.clearingStartTime !== undefined) {
      obj.clearingStartTime = message.clearingStartTime.toISOString();
    }
    if (message.clearingEndTime !== undefined) {
      obj.clearingEndTime = message.clearingEndTime.toISOString();
    }
    if (message.premarketStartTime !== undefined) {
      obj.premarketStartTime = message.premarketStartTime.toISOString();
    }
    if (message.premarketEndTime !== undefined) {
      obj.premarketEndTime = message.premarketEndTime.toISOString();
    }
    if (message.closingAuctionStartTime !== undefined) {
      obj.closingAuctionStartTime = message.closingAuctionStartTime.toISOString();
    }
    if (message.peningAuctionEndTime !== undefined) {
      obj.peningAuctionEndTime = message.peningAuctionEndTime.toISOString();
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TradingDay>, I>>(base?: I): TradingDay {
    return TradingDay.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TradingDay>, I>>(object: I): TradingDay {
    const message = createBaseTradingDay();
    message.date = object.date ?? undefined;
    message.isTradingDay = object.isTradingDay ?? false;
    message.startTime = object.startTime ?? undefined;
    message.endTime = object.endTime ?? undefined;
    message.openingAuctionStartTime = object.openingAuctionStartTime ?? undefined;
    message.closingAuctionEndTime = object.closingAuctionEndTime ?? undefined;
    message.eveningOpeningAuctionStartTime = object.eveningOpeningAuctionStartTime ?? undefined;
    message.eveningStartTime = object.eveningStartTime ?? undefined;
    message.eveningEndTime = object.eveningEndTime ?? undefined;
    message.clearingStartTime = object.clearingStartTime ?? undefined;
    message.clearingEndTime = object.clearingEndTime ?? undefined;
    message.premarketStartTime = object.premarketStartTime ?? undefined;
    message.premarketEndTime = object.premarketEndTime ?? undefined;
    message.closingAuctionStartTime = object.closingAuctionStartTime ?? undefined;
    message.peningAuctionEndTime = object.peningAuctionEndTime ?? undefined;
    return message;
  },
};

export type SharesService = typeof SharesService;
export const SharesService = {
  getShares: {
    path: "/shares.Shares/GetShares",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: GetInstrumentsRequest) => Buffer.from(GetInstrumentsRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => GetInstrumentsRequest.decode(value),
    responseSerialize: (value: GetSharesResponse) => Buffer.from(GetSharesResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => GetSharesResponse.decode(value),
  },
  getTradingSchedules: {
    path: "/shares.Shares/GetTradingSchedules",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: GetTradingSchedulesRequest) =>
      Buffer.from(GetTradingSchedulesRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => GetTradingSchedulesRequest.decode(value),
    responseSerialize: (value: GetTradingSchedulesResponse) =>
      Buffer.from(GetTradingSchedulesResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => GetTradingSchedulesResponse.decode(value),
  },
} as const;

export interface SharesServer extends UntypedServiceImplementation {
  getShares: handleUnaryCall<GetInstrumentsRequest, GetSharesResponse>;
  getTradingSchedules: handleUnaryCall<GetTradingSchedulesRequest, GetTradingSchedulesResponse>;
}

export interface SharesClient extends Client {
  getShares(
    request: GetInstrumentsRequest,
    callback: (error: ServiceError | null, response: GetSharesResponse) => void,
  ): ClientUnaryCall;
  getShares(
    request: GetInstrumentsRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: GetSharesResponse) => void,
  ): ClientUnaryCall;
  getShares(
    request: GetInstrumentsRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: GetSharesResponse) => void,
  ): ClientUnaryCall;
  getTradingSchedules(
    request: GetTradingSchedulesRequest,
    callback: (error: ServiceError | null, response: GetTradingSchedulesResponse) => void,
  ): ClientUnaryCall;
  getTradingSchedules(
    request: GetTradingSchedulesRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: GetTradingSchedulesResponse) => void,
  ): ClientUnaryCall;
  getTradingSchedules(
    request: GetTradingSchedulesRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: GetTradingSchedulesResponse) => void,
  ): ClientUnaryCall;
}

export const SharesClient = makeGenericClientConstructor(SharesService, "shares.Shares") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): SharesClient;
  service: typeof SharesService;
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
