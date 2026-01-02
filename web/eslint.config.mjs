import nextCoreWebVitals from "eslint-config-next/core-web-vitals";
import nextTypescript from "eslint-config-next/typescript";
import path from 'path';
import { fileURLToPath } from 'url';
import simpleImportSort from 'eslint-plugin-simple-import-sort';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const eslintConfig = [...nextCoreWebVitals, ...nextTypescript, {
  plugins: {
    'simple-import-sort': simpleImportSort,
  },
  rules: {
    'simple-import-sort/imports': [
      'error',
      {
        groups: [
          ['^react', '^next', '^@?\\w'],
          ['^@/app'],
          ['^@/core'],
          ['^@/features'],
          ['^$'],
          ['^@/shared'],
          ['^'],
        ],
      },
    ],
    'simple-import-sort/exports': 'error',
    '@typescript-eslint/consistent-type-definitions': ['error', 'type'],
    '@typescript-eslint/consistent-type-imports': [
      'error',
      { prefer: 'type-imports' },
    ],
  },
}, {
  ignores: ['.next/**', 'out/**', 'build/**', 'next-env.d.ts'],
}];

export default eslintConfig;
