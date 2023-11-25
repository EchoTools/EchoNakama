// characterizationUtils.ts

import fs from 'fs';
import path from 'path';

export const captureOutput = () => {
  // Simulate the request and capture the system output
  // Replace this with the actual logic to simulate the request in your system
  const simulatedRequest = /* replace with the logic to simulate the request */;
  return simulatedRequest;
};

export const compareWithGolden = (featureName: string, systemOutput: string) => {
  const goldenPath = path.join(__dirname, `../../golden/${featureName}.golden.txt`);
  const goldenContent = fs.readFileSync(goldenPath, 'utf-8').trim();

  // Compare the system output with the golden master
  if (systemOutput !== goldenContent) {
    // Optionally, update the golden master if needed
    fs.writeFileSync(goldenPath, systemOutput, 'utf-8');
  }

  return systemOutput;
};
