-- CreateTable
CREATE TABLE "Payment" (
    "id" TEXT NOT NULL,
    "create_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,
    "student_id" TEXT NOT NULL,
    "Amount" DOUBLE PRECISION NOT NULL,
    "Year" TEXT NOT NULL,
    "Semester" INTEGER NOT NULL,
    "Paid" BOOLEAN NOT NULL DEFAULT false,
    "payment_date" TIMESTAMP(3),

    CONSTRAINT "Payment_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Payment_student_id_Semester_Year_key" ON "Payment"("student_id", "Semester", "Year");

-- AddForeignKey
ALTER TABLE "Payment" ADD CONSTRAINT "Payment_student_id_fkey" FOREIGN KEY ("student_id") REFERENCES "Student"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
