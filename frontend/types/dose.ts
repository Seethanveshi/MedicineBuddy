export type Dose = {
    id : string;
    medicine_id : string;
    scheduled_at : string;
    status : "pending" | "taken" | "missed" | "skipped";
    taken_at : string | null;
}