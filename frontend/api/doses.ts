import { Dose } from "@/types/dose";
import axios from "axios";
import { format } from "date-fns";

export type CreateMedicineRequest = {
  name: string;
  dosage: string;
  start_date: string;
  end_date: string | null;
  schedule: {
    time: string;
    days_of_week: number[];
  };
};

export const API = axios.create({
    baseURL: "http://10.81.193.101:8080",
});

export const getTodayDoses = async (): Promise<Dose[]> => {
    const res = await API.get<Dose[]>("/doses/today");
    return res.data;
}

export const takeDose = async (id: string) => {
    await API.post(`doses/${id}/take`);   
}

export const skipDose = async (id: string) => {
    await API.post(`/doses/${id}/skip`);
}

export const getUpcomingDoses = async (days = 7): Promise<Dose[]> => {
  const res = await API.get<Dose[]>(`/doses/upcoming?days=${days}`);
  return res.data;
};

export const getDosesByDate = async (date: Date): Promise<Dose[]> => {
  const formatted = format(date, "yyyy-MM-dd");

  const res = await API.get<Dose[]>(`/doses/date?date=${formatted}`);
  return res.data;
};


export const createMedicine = async (
  data: CreateMedicineRequest
) => {
  try {
    console.log(data);
    const res = await API.post("/medicines", data);
    return res.data;
  }  catch (error) {
    console.error("error while sending request");
  }
};