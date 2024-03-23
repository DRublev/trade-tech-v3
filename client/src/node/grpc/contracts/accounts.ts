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

export const protobufPackage = "accounts";

export interface Account {
  Id: string;
  Name: string;
}

export interface GetAccountsRequest {
}

export interface GetAccountsResponse {
  Accounts: Account[];
}

export interface SetAccountRequest {
  AccountId: string;
}

export interface SetAccountResponse {
}

function createBaseAccount(): Account {
  return { Id: "", Name: "" };
}

export const Account = {
  encode(message: Account, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.Id !== "") {
      writer.uint32(10).string(message.Id);
    }
    if (message.Name !== "") {
      writer.uint32(18).string(message.Name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Account {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAccount();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Id = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.Name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Account {
    return {
      Id: isSet(object.Id) ? globalThis.String(object.Id) : "",
      Name: isSet(object.Name) ? globalThis.String(object.Name) : "",
    };
  },

  toJSON(message: Account): unknown {
    const obj: any = {};
    if (message.Id !== "") {
      obj.Id = message.Id;
    }
    if (message.Name !== "") {
      obj.Name = message.Name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Account>, I>>(base?: I): Account {
    return Account.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Account>, I>>(object: I): Account {
    const message = createBaseAccount();
    message.Id = object.Id ?? "";
    message.Name = object.Name ?? "";
    return message;
  },
};

function createBaseGetAccountsRequest(): GetAccountsRequest {
  return {};
}

export const GetAccountsRequest = {
  encode(_: GetAccountsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAccountsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAccountsRequest();
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

  fromJSON(_: any): GetAccountsRequest {
    return {};
  },

  toJSON(_: GetAccountsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAccountsRequest>, I>>(base?: I): GetAccountsRequest {
    return GetAccountsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAccountsRequest>, I>>(_: I): GetAccountsRequest {
    const message = createBaseGetAccountsRequest();
    return message;
  },
};

function createBaseGetAccountsResponse(): GetAccountsResponse {
  return { Accounts: [] };
}

export const GetAccountsResponse = {
  encode(message: GetAccountsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.Accounts) {
      Account.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAccountsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAccountsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.Accounts.push(Account.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAccountsResponse {
    return {
      Accounts: globalThis.Array.isArray(object?.Accounts) ? object.Accounts.map((e: any) => Account.fromJSON(e)) : [],
    };
  },

  toJSON(message: GetAccountsResponse): unknown {
    const obj: any = {};
    if (message.Accounts?.length) {
      obj.Accounts = message.Accounts.map((e) => Account.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAccountsResponse>, I>>(base?: I): GetAccountsResponse {
    return GetAccountsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAccountsResponse>, I>>(object: I): GetAccountsResponse {
    const message = createBaseGetAccountsResponse();
    message.Accounts = object.Accounts?.map((e) => Account.fromPartial(e)) || [];
    return message;
  },
};

function createBaseSetAccountRequest(): SetAccountRequest {
  return { AccountId: "" };
}

export const SetAccountRequest = {
  encode(message: SetAccountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.AccountId !== "") {
      writer.uint32(10).string(message.AccountId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetAccountRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetAccountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.AccountId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetAccountRequest {
    return { AccountId: isSet(object.AccountId) ? globalThis.String(object.AccountId) : "" };
  },

  toJSON(message: SetAccountRequest): unknown {
    const obj: any = {};
    if (message.AccountId !== "") {
      obj.AccountId = message.AccountId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetAccountRequest>, I>>(base?: I): SetAccountRequest {
    return SetAccountRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetAccountRequest>, I>>(object: I): SetAccountRequest {
    const message = createBaseSetAccountRequest();
    message.AccountId = object.AccountId ?? "";
    return message;
  },
};

function createBaseSetAccountResponse(): SetAccountResponse {
  return {};
}

export const SetAccountResponse = {
  encode(_: SetAccountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetAccountResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetAccountResponse();
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

  fromJSON(_: any): SetAccountResponse {
    return {};
  },

  toJSON(_: SetAccountResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<SetAccountResponse>, I>>(base?: I): SetAccountResponse {
    return SetAccountResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetAccountResponse>, I>>(_: I): SetAccountResponse {
    const message = createBaseSetAccountResponse();
    return message;
  },
};

export type AccountsService = typeof AccountsService;
export const AccountsService = {
  getAccounts: {
    path: "/accounts.Accounts/GetAccounts",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: GetAccountsRequest) => Buffer.from(GetAccountsRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => GetAccountsRequest.decode(value),
    responseSerialize: (value: GetAccountsResponse) => Buffer.from(GetAccountsResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => GetAccountsResponse.decode(value),
  },
  setAccount: {
    path: "/accounts.Accounts/SetAccount",
    requestStream: false,
    responseStream: false,
    requestSerialize: (value: SetAccountRequest) => Buffer.from(SetAccountRequest.encode(value).finish()),
    requestDeserialize: (value: Buffer) => SetAccountRequest.decode(value),
    responseSerialize: (value: SetAccountResponse) => Buffer.from(SetAccountResponse.encode(value).finish()),
    responseDeserialize: (value: Buffer) => SetAccountResponse.decode(value),
  },
} as const;

export interface AccountsServer extends UntypedServiceImplementation {
  getAccounts: handleUnaryCall<GetAccountsRequest, GetAccountsResponse>;
  setAccount: handleUnaryCall<SetAccountRequest, SetAccountResponse>;
}

export interface AccountsClient extends Client {
  getAccounts(
    request: GetAccountsRequest,
    callback: (error: ServiceError | null, response: GetAccountsResponse) => void,
  ): ClientUnaryCall;
  getAccounts(
    request: GetAccountsRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: GetAccountsResponse) => void,
  ): ClientUnaryCall;
  getAccounts(
    request: GetAccountsRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: GetAccountsResponse) => void,
  ): ClientUnaryCall;
  setAccount(
    request: SetAccountRequest,
    callback: (error: ServiceError | null, response: SetAccountResponse) => void,
  ): ClientUnaryCall;
  setAccount(
    request: SetAccountRequest,
    metadata: Metadata,
    callback: (error: ServiceError | null, response: SetAccountResponse) => void,
  ): ClientUnaryCall;
  setAccount(
    request: SetAccountRequest,
    metadata: Metadata,
    options: Partial<CallOptions>,
    callback: (error: ServiceError | null, response: SetAccountResponse) => void,
  ): ClientUnaryCall;
}

export const AccountsClient = makeGenericClientConstructor(AccountsService, "accounts.Accounts") as unknown as {
  new (address: string, credentials: ChannelCredentials, options?: Partial<ClientOptions>): AccountsClient;
  service: typeof AccountsService;
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
