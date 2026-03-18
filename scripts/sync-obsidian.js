#!/usr/bin/env node

// Auto-generated stub for Obsidian vault synchronization
// Syncs notes from resources/notes/ to Obsidian workspace

import { readdir, readFile, writeFile, mkdir, stat } from 'fs/promises';
import { join, relative, extname, basename, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const SOURCE_DIR = join(process.cwd(), 'resources/notes');
const TARGET_DIR = process.env.OBSIDIAN_VAULT_PATH || join(process.cwd(), 'obsidian-vault');

async function ensureDir(dirPath) {
  try {
    await mkdir(dirPath, { recursive: true });
  } catch (err) {
    if (err.code !== 'EEXIST') throw err;
  }
}

async function syncFile(srcPath, targetPath) {
  const content = await readFile(srcPath, 'utf-8');
  await ensureDir(dirname(targetPath));
  await writeFile(targetPath, content);
  console.log(`Synced: ${relative(process.cwd(), srcPath)} -> ${relative(process.cwd(), targetPath)}`);
}

async function syncRecursive(srcDir, targetDir) {
  const entries = await readdir(srcDir, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = join(srcDir, entry.name);
    const relPath = relative(SOURCE_DIR, srcPath);
    const targetPath = join(TARGET_DIR, relPath);

    if (entry.isDirectory()) {
      await syncRecursive(srcPath, targetPath);
    } else if (entry.isFile()) {
      await syncFile(srcPath, targetPath);
    }
  }
}

async function run() {
  console.log('Starting Obsidian sync...');
  try {
    await syncRecursive(SOURCE_DIR, TARGET_DIR);
    console.log('Obsidian sync completed successfully!');
  } catch (error) {
    console.error('Obsidian sync failed:', error);
    process.exit(1);
  }
}

run();