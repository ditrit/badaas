/* eslint-disable no-console */

const path = require('path');
const { promisify } = require('util');

const render = require('@koa/ejs');
const helmet = require('helmet');

const { Provider } = require('oidc-provider');

const configuration = require('./support/configuration');
const routes = require('./routes/koa');

const PORT = process.env.PORT || 3000;
const ISSUER = `http://${process.env.HOST || "localhost"}:${PORT}`;
const Account = require('./support/account');
const { Console } = require('console');
configuration.findAccount = Account.findAccount;
configuration.pkce = {
  required: function pkceRequired(ctx, client) {
    return false;
  }
}

let server;

console.log(`oidc-provider conf`);
console.log(`PORT: ${PORT}`);
console.log(`ISSUER: ${ISSUER}`);
console.log(`CLIENT_ID: ${process.env.CLIENT_ID || 'oidcCLIENT'}`);
console.log(`CLIENT_SECRET: ${process.env.CLIENT_SECRET ||'abcd'}`);
console.log(`REDIRECT_URI: ${process.env.REDIRECT_URI ||'http://localhost:8002/auth/oidc/callback'}`);

(async () => {
  let adapter;
  if (process.env.MONGODB_URI) {
    adapter = require('./adapters/mongodb'); // eslint-disable-line global-require
    await adapter.connect();
  }

  const prod = process.env.NODE_ENV === 'production';

  const provider = new Provider(ISSUER, { adapter, ...configuration });

  const directives = helmet.contentSecurityPolicy.getDefaultDirectives();
  delete directives['form-action'];
  const pHelmet = promisify(helmet({
    contentSecurityPolicy: {
      useDefaults: false,
      directives,
    },
  }));

  provider.use(async (ctx, next) => {
    const origSecure = ctx.req.secure;
    ctx.req.secure = ctx.request.secure;
    await pHelmet(ctx.req, ctx.res);
    ctx.req.secure = origSecure;
    return next();
  });

  if (prod) {
    provider.proxy = true;
    provider.use(async (ctx, next) => {
      if (ctx.secure) {
        await next();
      } else if (ctx.method === 'GET' || ctx.method === 'HEAD') {
        ctx.status = 303;
        ctx.redirect(ctx.href.replace(/^http:\/\//i, 'https://'));
      } else {
        ctx.body = {
          error: 'invalid_request',
          error_description: 'do yourself a favor and only use https',
        };
        ctx.status = 400;
      }
    });
  }
  render(provider.app, {
    cache: false,
    viewExt: 'ejs',
    layout: '_layout',
    root: path.join(__dirname, 'views'),
  });
  provider.use(routes(provider).routes());
  server = provider.listen(PORT, () => {
    console.log(`application is listening on port ${PORT}, check its /.well-known/openid-configuration`);
  });
})().catch((err) => {
  if (server && server.listening) server.close();
  console.error(err);
  process.exitCode = 1;
});
