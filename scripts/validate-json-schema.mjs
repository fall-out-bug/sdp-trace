#!/usr/bin/env node
import fs from 'node:fs';
import path from 'node:path';
import process from 'node:process';
import Ajv2020 from 'ajv/dist/2020.js';

const [schemaPath, dataPath] = process.argv.slice(2);

if (!schemaPath || !dataPath) {
  console.error('Usage: validate-json-schema.mjs <schema> <data>');
  process.exit(2);
}

const readJson = (file) => JSON.parse(fs.readFileSync(file, 'utf8'));
const schema = readJson(schemaPath);
const data = readJson(dataPath);

const ajv = new Ajv2020({
  allErrors: true,
  strict: false,
  validateSchema: true
});

for (const file of fs.readdirSync('schema').filter((name) => name.endsWith('.json')).sort()) {
  const refPath = path.join('schema', file);
  const refSchema = readJson(refPath);
  if (refPath !== schemaPath && refSchema.$id) {
    ajv.addSchema(refSchema);
  }
}

const validate = ajv.compile(schema);
if (!validate(data)) {
  console.error(`${dataPath} invalid`);
  console.error(JSON.stringify(validate.errors, null, 2));
  process.exit(1);
}

console.log(`${dataPath} valid`);
