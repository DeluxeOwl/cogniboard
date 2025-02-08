import { z } from 'zod'

export const errorDetailSchema = z.object({
  location: z.string().describe("Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id'").optional(),
  message: z.string().describe('Error message text').optional(),
  value: z.any().optional(),
})