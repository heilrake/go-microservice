export type StaticRoute<T extends string> = () => T;

export type ParamRoute<
  TPrefix extends string,
  TParam
> = (param: TParam) => `${TPrefix}/${string}`;