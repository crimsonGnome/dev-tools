# Frontend Package Steering

This package is the Vigilance React Native app (branding: Vigilon/Ignisight).

## Stack
- React Native 0.74.3, TypeScript
- styled-components for styling
- React Navigation (native-stack)

## Rules
- All components must be typed — no `any`
- Use styled-components for all styling; no inline style objects
- Navigation changes require updating the type definitions in the navigator
- Test on both iOS and Android before marking a feature complete
