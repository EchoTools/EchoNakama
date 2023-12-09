import 'jest-ts-auto-mock';

module.exports = {
    transform: {
      "^.+\\.tsx?$": "ts-jest",
    },
    testMatch: ['<rootDir>/src/**/*.test.ts'],
    testRegex: "(/__tests__/.*|(\\.|/)(test|spec))\\.(jsx?|tsx?)$",
    moduleFileExtensions: ["ts", "tsx", "js", "jsx", "json", "node"],
    collectCoverage: true,
    mapCoverage: true,
  };