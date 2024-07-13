import { record, z } from "zod";

// We're keeping a simple non-relational schema here.
// IRL, you will have a schema for your data models.
export const taskSchema = z.object({
  id: z.string(),
  title: z.string(),
  status: z.string(),
  label: z.string(),
  priority: z.string(),
});

export const domainSchema = z.object({
  id: z.string(),
  type: z.string(),
  domain: z.string(),
  record: z.string(),
  status: z.number(),
  isp: z.string(),
  ip: z.string(),
  city: z.string(),
  expire_at: z.string(),
});

export type Task = z.infer<typeof taskSchema>;
export type Domain = z.infer<typeof domainSchema>;
