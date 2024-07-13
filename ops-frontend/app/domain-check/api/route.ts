import { NextResponse } from "next/server";

import { Domains } from "@/app/domain-check/api/domain";

export async function GET() {
  return NextResponse.json(Domains, { status: 200 });
}
