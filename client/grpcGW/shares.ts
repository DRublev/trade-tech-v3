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

export const protobufPackage = "shares";

export interface Share {
  Figi: string;
  Name: string;
  /** ticker, lot, ipo_date, trading_status, min_price_increment, uid, first_1min_candle_date, first_1day_candle_date */
  Exchange: string;
}

export interface GetInstrumentsRequest {
  instrumentStatus: number;
}

export interface GetSharesResponse {
  instruments: Share[];
}

function createBaseShare(): Share {
  return { Figi: "", Name: "", Exchange: "" };
}

export const Share = {
  encode(message: Share, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Figi !== "") {
      writer.uint32(10).string(message.Figi);
    }
    if (message.Name !== "") {
      writer.uint32(18).string(message.Name);
    }
    if (message.Exchange !== "") {
      writer.uint32(26).string(message.Exchange);
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

          message.Figi = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.Exchange = reader.string();
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
      Figi: isSet(object.Figi) ? globalThis.String(object.Figi) : "",
      Name: isSet(object.Name) ? globalThis.String(object.Name) : "",
      Exchange: isSet(object.Exchange) ? globalThis.String(object.Exchange) : "",
    };
  },

  toJSON(message: Share): unknown {
    const obj: any = {};
    if (message.Figi !== "") {
      obj.Figi = message.Figi;
    }
    if (message.Name !== "") {
      obj.Name = message.Name;
    }
    if (message.Exchange !== "") {
      obj.Exchange = message.Exchange;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Share>, I>>(base?: I): Share {
    return Share.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Share>, I>>(object: I): Share {
    const message = createBaseShare();
    message.Figi = object.Figi ?? "";
    message.Name = object.Name ?? "";
    message.Exchange = object.Exchange ?? "";
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
} as const;

export interface SharesServer extends UntypedServiceImplementation {
  getShares: handleUnaryCall<GetInstrumentsRequest, GetSharesResponse>;
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
