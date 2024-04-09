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

export const protobufPackage = "auth";

export interface SetTokenRequest {
  Token: string;
  IsSandbox: boolean;
}

export interface SetTokenResponse {
}

export interface ClearTokenRequest {
  ForSandbox: boolean;
}

export interface ClearTokenResponse {
}

export interface HasTokenRequest {
}

export interface HasTokenResponse {
  HasToken: boolean;
}

export interface PruneTokensRequest {
}

export interface PruneTokensResponse {
}

function createBaseSetTokenRequest(): SetTokenRequest {
  return { Token: "", IsSandbox: false };
}

export const SetTokenRequest = {
  encode(message: SetTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Token !== "") {
      writer.uint32(10).string(message.Token);
    }
    if (message.IsSandbox !== false) {
      writer.uint32(16).bool(message.IsSandbox);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Token = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.IsSandbox = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetTokenRequest {
    return {
      Token: isSet(object.Token) ? globalThis.String(object.Token) : "",
      IsSandbox: isSet(object.IsSandbox) ? globalThis.Boolean(object.IsSandbox) : false,
    };
  },

  toJSON(message: SetTokenRequest): unknown {
    const obj: any = {};
    if (message.Token !== "") {
      obj.Token = message.Token;
    }
    if (message.IsSandbox !== false) {
      obj.IsSandbox = message.IsSandbox;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetTokenRequest>, I>>(base?: I): SetTokenRequest {
    return SetTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetTokenRequest>, I>>(object: I): SetTokenRequest {
    const message = createBaseSetTokenRequest();
    message.Token = object.Token ?? "";
    message.IsSandbox = object.IsSandbox ?? false;
    return message;
  },
};

function createBaseSetTokenResponse(): SetTokenResponse {
  return {};
}

export const SetTokenResponse = {
  encode(_: SetTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetTokenResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): SetTokenResponse {
    return {};
  },

  toJSON(_: SetTokenResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<SetTokenResponse>, I>>(base?: I): SetTokenResponse {
    return SetTokenResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetTokenResponse>, I>>(_: I): SetTokenResponse {
    const message = createBaseSetTokenResponse();
    return message;
  },
};

function createBaseClearTokenRequest(): ClearTokenRequest {
  return { ForSandbox: false };
}

export const ClearTokenRequest = {
  encode(message: ClearTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ForSandbox !== false) {
      writer.uint32(8).bool(message.ForSandbox);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClearTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClearTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.ForSandbox = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ClearTokenRequest {
    return { ForSandbox: isSet(object.ForSandbox) ? globalThis.Boolean(object.ForSandbox) : false };
  },

  toJSON(message: ClearTokenRequest): unknown {
    const obj: any = {};
    if (message.ForSandbox !== false) {
      obj.ForSandbox = message.ForSandbox;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ClearTokenRequest>, I>>(base?: I): ClearTokenRequest {
    return ClearTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ClearTokenRequest>, I>>(object: I): ClearTokenRequest {
    const message = createBaseClearTokenRequest();
    message.ForSandbox = object.ForSandbox ?? false;
    return message;
  },
};

function createBaseClearTokenResponse(): ClearTokenResponse {
  return {};
}

export const ClearTokenResponse = {
  encode(_: ClearTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClearTokenResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClearTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ClearTokenResponse {
    return {};
  },

  toJSON(_: ClearTokenResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<ClearTokenResponse>, I>>(base?: I): ClearTokenResponse {
    return ClearTokenResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ClearTokenResponse>, I>>(_: I): ClearTokenResponse {
    const message = createBaseClearTokenResponse();
    return message;
  },
};

function createBaseHasTokenRequest(): HasTokenRequest {
  return {};
}

export const HasTokenRequest = {
  encode(_: HasTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HasTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHasTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): HasTokenRequest {
    return {};
  },

  toJSON(_: HasTokenRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<HasTokenRequest>, I>>(base?: I): HasTokenRequest {
    return HasTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HasTokenRequest>, I>>(_: I): HasTokenRequest {
    const message = createBaseHasTokenRequest();
    return message;
  },
};

function createBaseHasTokenResponse(): HasTokenResponse {
  return { HasToken: false };
}

export const HasTokenResponse = {
  encode(message: HasTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.HasToken !== false) {
      writer.uint32(8).bool(message.HasToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HasTokenResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHasTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.HasToken = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HasTokenResponse {
    return { HasToken: isSet(object.HasToken) ? globalThis.Boolean(object.HasToken) : false };
  },

  toJSON(message: HasTokenResponse): unknown {
    const obj: any = {};
    if (message.HasToken !== false) {
      obj.HasToken = message.HasToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HasTokenResponse>, I>>(base?: I): HasTokenResponse {
    return HasTokenResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HasTokenResponse>, I>>(object: I): HasTokenResponse {
    const message = createBaseHasTokenResponse();
    message.HasToken = object.HasToken ?? false;
    return message;
  },
};

function createBasePruneTokensRequest(): PruneTokensRequest {
  return {};
}

export const PruneTokensRequest = {
  encode(_: PruneTokensRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PruneTokensRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePruneTokensRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PruneTokensRequest {
    return {};
  },

  toJSON(_: PruneTokensRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PruneTokensRequest>, I>>(base?: I): PruneTokensRequest {
    return PruneTokensRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PruneTokensRequest>, I>>(_: I): PruneTokensRequest {
    const message = createBasePruneTokensRequest();
    return message;
  },
};

function createBasePruneTokensResponse(): PruneTokensResponse {
  return {};
}

export const PruneTokensResponse = {
  encode(_: PruneTokensResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PruneTokensResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePruneTokensResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PruneTokensResponse {
    return {};
  },

  toJSON(_: PruneTokensResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PruneTokensResponse>, I>>(base?: I): PruneTokensResponse {
    return PruneTokensResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<PruneTokensResponse>, I>>(_: I): PruneTokensResponse {
    const message = createBasePruneTokensResponse();
    return message;
  },
};

export type AuthService = typeof AuthService;
export const AuthService = {
  setToken: {
    path: "/auth.Auth/SetToken",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: SetTokenRequest) => Buffer.from(SetTokenRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SetTokenRequest.decode(value),
    responseSerialize: (value: SetTokenResponse) => Buffer.from(SetTokenResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => SetTokenResponse.decode(value),
  },
  clearToken: {
    path: "/auth.Auth/ClearToken",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: ClearTokenRequest) => Buffer.from(ClearTokenRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => ClearTokenRequest.decode(value),
    responseSerialize: (value: ClearTokenResponse) => Buffer.from(ClearTokenResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => ClearTokenResponse.decode(value),
  },
  hasToken: {
    path: "/auth.Auth/HasToken",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: HasTokenRequest) => Buffer.from(HasTokenRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => HasTokenRequest.decode(value),
    responseSerialize: (value: HasTokenResponse) => Buffer.from(HasTokenResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => HasTokenResponse.decode(value),
  },
  pruneTokens: {
    path: "/auth.Auth/PruneTokens",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: PruneTokensRequest) => Buffer.from(PruneTokensRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => PruneTokensRequest.decode(value),
    responseSerialize: (value: PruneTokensResponse) => Buffer.from(PruneTokensResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => PruneTokensResponse.decode(value),
  },
} as const;

export interface AuthServer extends UntypedServiceImplementation {
  setToken: handleUnaryCall<SetTokenRequest, SetTokenResponse>;
  clearToken: handleUnaryCall<ClearTokenRequest, ClearTokenResponse>;
  hasToken: handleUnaryCall<HasTokenRequest, HasTokenResponse>;
  pruneTokens: handleUnaryCall<PruneTokensRequest, PruneTokensResponse>;
}

export interface AuthClient extends Client {
  setToken(
    request: SetTokenRequest,
    callback: (error: ServiceError | null, response: SetTokenResponse) => void,
  ): ClientUnaryCall;
  setToken(
    request: SetTokenRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: SetTokenResponse) => void,
  ): ClientUnaryCall;
  setToken(
    request: SetTokenRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: SetTokenResponse) => void,
  ): ClientUnaryCall;
  clearToken(
    request: ClearTokenRequest,
    callback: (error: ServiceError | null, response: ClearTokenResponse) => void,
  ): ClientUnaryCall;
  clearToken(
    request: ClearTokenRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: ClearTokenResponse) => void,
  ): ClientUnaryCall;
  clearToken(
    request: ClearTokenRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: ClearTokenResponse) => void,
  ): ClientUnaryCall;
  hasToken(
    request: HasTokenRequest,
    callback: (error: ServiceError | null, response: HasTokenResponse) => void,
  ): ClientUnaryCall;
  hasToken(
    request: HasTokenRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: HasTokenResponse) => void,
  ): ClientUnaryCall;
  hasToken(
    request: HasTokenRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: HasTokenResponse) => void,
  ): ClientUnaryCall;
  pruneTokens(
    request: PruneTokensRequest,
    callback: (error: ServiceError | null, response: PruneTokensResponse) => void,
  ): ClientUnaryCall;
  pruneTokens(
    request: PruneTokensRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: PruneTokensResponse) => void,
  ): ClientUnaryCall;
  pruneTokens(
    request: PruneTokensRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: PruneTokensResponse) => void,
  ): ClientUnaryCall;
}

export const AuthClient = makeGenericClientConstructor(AuthService, "auth.Auth") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): AuthClient;
  service: typeof AuthService;
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
