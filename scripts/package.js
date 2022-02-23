const { join } = require('path');
const { cwd } = require('process');
const { path7za } = require('7zip-bin');
const { spawnSync, execSync } = require('child_process');
const { mkdirSync, mkdtempSync } = require('fs');
const os = require('os');
const fse = require('fs-extra');
const rimraf = require('rimraf');
const webappOptions = require('../apps/webapp/webpack.config.js');
const webpack = require('webpack');

const distPath = join(cwd(), 'dist');
const sfxPath = join(cwd(), 'sfx');
const serverPath = join(cwd(), 'apps', 'server', 'main');

rimraf.sync(distPath);

const run = (compiler) =>
  new Promise((resolve) => {
    compiler.run(() => {
      resolve();
    });
  });

const tmpPathWin64 = mkdtempSync(join(os.tmpdir(), 'build-win64-'));
const tmpPathLinux64 = mkdtempSync(join(os.tmpdir(), 'build-linux64-'));
fse.copySync(tmpPathLinux64, distPath);

(async () => {
  //Build webapp
  await Promise.all([run(webpack(webappOptions))]);

  //Build server for windows

  execSync(
    `cd ${serverPath} && cross-env GOOS=windows GOARCH=amd64 go build -ldflags=\"-s -w\" -o ${tmpPathWin64}/main.exe`,
    {
      stdio: 'inherit',
    },
  );

  //Package webapp+server for windows

  //Make 7z archive
  spawnSync(
    path7za,
    [
      'a',
      //Output
      `${tmpPathWin64}\\package-win64.7z`,
      //Server file
      `${tmpPathWin64}\\main.exe`,
      //Webapp files
      `${distPath}\\*`,
    ],
    {
      stdio: 'inherit',
    },
  );

  //Make sfx exe
  execSync(
    `COPY /b "${sfxPath}\\7z.sfx" + "${sfxPath}\\sfx.txt" + "${tmpPathWin64}\\package-win64.7z" "${distPath}\\package-win64.exe"`,
    {
      stdio: 'inherit',
    },
  );

  //Build server for linux
  execSync(
    `cd ${serverPath} && cross-env GOOS=linux GOARCH=amd64 go build -ldflags=\"-s -w\" -o ${tmpPathLinux64}/main`,
    {
      stdio: 'inherit',
    },
  );

  //Package for linux
  execSync(
    `${sfxPath}\\makeself\\makeself.sh --notemp ${tmpPathLinux64} ${distPath}\\package-linux64 "package-linux64" ./main`,
    {
      stdio: 'inherit',
    },
  );

  rimraf.sync(tmpPathLinux64);
  rimraf.sync(tmpPathWin64);
})();
