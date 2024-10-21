import bodyParser from 'body-parser';
import cors from 'cors';
import express, { Express } from 'express';
import expressRequestId from 'express-request-id';
import nocache from 'nocache';
import responseTime from 'response-time';

import RequestLogger from './middlewares/request-logger';
import resBodyCapture from './middlewares/res-body-capture';

export default function server(): Express {
  const app = express();

  app.use(responseTime({ suffix: false }));

  app.use(expressRequestId());

  app.use(resBodyCapture);

  const corsOptions = {
    // origin: 'https://testnet-dapp.nemonetwork.org',
    origin: '*', // Allow all origins
    methods: ['GET', 'HEAD', 'PUT', 'PATCH', 'POST', 'DELETE', 'OPTIONS'], // Allow all methods
    allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With', 'Origin', 'Accept'], // Allow all common headers
    exposedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With', 'Origin', 'Accept'], // Expose headers to the browser
    credentials: true, // Allow credentials (cookies, authorization headers)
    preflightContinue: true, // Pass the CORS preflight response to the next handler
    optionsSuccessStatus: 204,
  };

  app.use(cors(corsOptions));
  // Handle preflight OPTIONS requests
  app.options('*', cors(corsOptions));
  app.use(nocache());

  app.get('/health', (_req: express.Request, res: express.Response) => {
    return res.status(200).json({ ok: true });
  });

  app.use((_req: express.Request, _res: express.Response, next: express.NextFunction) => next());

  app.use(bodyParser.json());

  app.use(RequestLogger);

  app.use((_req: express.Request, res: express.Response) => {
    res.status(404).json({
      error: 'Not Found',
      errorCode: 404,
    });
  });

  return app;
}
