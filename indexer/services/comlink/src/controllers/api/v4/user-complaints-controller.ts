import { logger, stats } from '@nemo-network-indexer/base/build';
import {
  UserComplaintFromDatabase,
  UserComplaintTable,
} from '@nemo-network-indexer/postgres/build/src';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';

import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';

const router: express.Router = express.Router();
const controllerName: string = 'user-complaints-controller';

const UserComplaintSchema = checkSchema({
  message: {
    in: ['body'],
    isString: true,
    notEmpty: true,
    errorMessage: 'message is required and must be a non-empty string',
  },
  email: {
    in: ['body'],
    optional: true,
    isEmail: true,
    errorMessage: 'email must be a valid email address',
  },
});

router.post(
  '/',
  ...UserComplaintSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      message,
      email,
    }: {
      message: string,
      email?: string,
    } = matchedData(req) as { message: string, email?: string };

    try {
      const complaint: UserComplaintFromDatabase = await UserComplaintTable.create({
        message: message.trim(),
        email,
      });

      return res.status(201).send(complaint);
    } catch (error) {
      logger.error({
        at: 'UserComplaintsController POST /',
        message: 'Failed to create user complaint',
        error,
        body: JSON.stringify(req.body),
      });
      return handleControllerError(
        'UserComplaintsController POST /',
        'Failed to create user complaint',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.create_complaint.timing`,
        Date.now() - start,
      );
    }
  },
);

export default router;
