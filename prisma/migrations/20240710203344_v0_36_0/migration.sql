-- DropIndex
DROP INDEX "StepRun_tenantId_status_requeueAfter_createdAt_idx";

-- CreateIndex
CREATE INDEX "StepRun_tenantId_idx" ON "StepRun"("tenantId");

-- CreateIndex
CREATE INDEX "StepRun_workerId_idx" ON "StepRun"("workerId");

-- CreateIndex
CREATE INDEX "StepRun_createdAt_idx" ON "StepRun"("createdAt");

-- CreateIndex
CREATE INDEX "StepRun_jobRunId_status_tenantId_requeueAfter_idx" ON "StepRun"("jobRunId", "status", "tenantId", "requeueAfter");

-- CreateIndex
CREATE INDEX "StepRun_status_timeoutAt_tickerId_idx" ON "StepRun"("status", "timeoutAt", "tickerId");
