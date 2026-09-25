/* eslint-disable no-unused-vars */
import type {JSX as ReactJSX} from 'react';

// @types/react 19 dropped the global JSX namespace, but react-markdown 8 ships
// lib/complex-types.ts (type-checked as source, so skipLibCheck doesn't apply) that still uses it.
declare global {
    namespace JSX {
        type Element = ReactJSX.Element;
        type IntrinsicElements = ReactJSX.IntrinsicElements;
    }
}
