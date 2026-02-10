import { Dose } from "@/types/dose";

export const groupByDate = (doses: Dose[]) => {
    const map: Record<string, Dose[]> = {};

    doses.forEach((dose) => {
      const date = new Date(dose.scheduled_at).toISOString().split("T")[0];

      if (!map[date]) map[date] = [];
      map[date].push(dose);
    });

    return Object.entries(map).sort(
      ([a], [b]) => new Date(a).getTime() - new Date(b).getTime()
    );
  };

  export const groupByTime = (doses: Dose[]) => {
      const map: Record<string, Dose[]> = {};
  
      doses.forEach((dose) => {
        const time = new Date(dose.scheduled_at).toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        });
  
        if (!map[time]) map[time] = [];
        map[time].push(dose);
      });
  
      return Object.entries(map).sort(
        ([a], [b]) =>
          Date.parse(`1970-01-01 ${a}`) - Date.parse(`1970-01-01 ${b}`)
      );
  
    };