import { logger, stats } from '@nemo-network-indexer/base/build';
import {
  UserComplaintFromDatabase,
  UserComplaintTable,
} from '@nemo-network-indexer/postgres/build/src';
import express from 'express';
import { checkSchema, matchedData } from 'express-validator';
import {
  Body, Controller, Get, Post, Query, Route, SuccessResponse,
} from 'tsoa';

import { getReqRateLimiter } from '../../../caches/rate-limiters';
import config from '../../../config';
import { handleControllerError } from '../../../lib/helpers';
import { rateLimiterMiddleware } from '../../../lib/rate-limit';
import { CheckLimitSchema } from '../../../lib/validation/schemas';
import { handleValidationErrors } from '../../../request-helpers/error-handler';
import ExportResponseCodeStats from '../../../request-helpers/export-response-code-stats';
import { UserComplaintsResponse } from '../../../types';

const router: express.Router = express.Router();
const controllerName: string = 'user-complaints-controller';

interface UserComplaintRequest {
  walletAddress: string,
  message: string,
  email?: string,
}

const UserComplaintSchema = checkSchema({
  walletAddress: {
    in: ['body'],
    isString: true,
    notEmpty: true,
    errorMessage: 'walletAddress is required and must be a non-empty string',
  },
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

@Route('userComplaints')
class UserComplaintsController extends Controller {
  @Get('/')
  async getAll(
    @Query() limit?: number,
  ): Promise<UserComplaintsResponse> {
    const complaints: UserComplaintFromDatabase[] = await UserComplaintTable.findAll(
      { limit },
      [],
    );
    return { complaints };
  }

  @Post('/')
  @SuccessResponse('201', 'Created')
  async create(
    @Body() body: UserComplaintRequest,
  ): Promise<UserComplaintFromDatabase> {
    return UserComplaintTable.create({
      walletAddress: body.walletAddress,
      message: body.message.trim(),
      email: body.email,
    });
  }
}

router.get(
  '/',
  rateLimiterMiddleware(getReqRateLimiter),
  ...CheckLimitSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();
    const { limit } = matchedData(req) as { limit?: number };

    try {
      const controller: UserComplaintsController = new UserComplaintsController();
      const response: UserComplaintsResponse = await controller.getAll(limit);

      return res.send(response);
    } catch (error) {
      return handleControllerError(
        'UserComplaintsController GET /',
        'Failed to fetch user complaints',
        error,
        req,
        res,
      );
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.${controllerName}.get_complaints.timing`,
        Date.now() - start,
      );
    }
  },
);

router.post(
  '/',
  ...UserComplaintSchema,
  handleValidationErrors,
  ExportResponseCodeStats({ controllerName }),
  async (req: express.Request, res: express.Response) => {
    const start: number = Date.now();

    const {
      walletAddress,
      message,
      email,
    } = matchedData(req) as UserComplaintRequest;

    try {
      const controller: UserComplaintsController = new UserComplaintsController();
      const complaint: UserComplaintFromDatabase = await controller.create({
        walletAddress,
        message,
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
