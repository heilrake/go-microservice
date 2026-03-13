export type StaticKey<T extends string> = () => T;

export type ParamKey<
  TPrefix extends string,
  TParam
> = (param: TParam) => readonly [TPrefix, TParam] | null;