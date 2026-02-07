import { useEffect, useState } from "react";
import { View, FlatList, Text } from "react-native";
import { getTodayDoses, takeDose, skipDose, getUpcomingDoses } from "../api/doses";
import DoseCard from "../components/DoseCard";
import { Dose } from "../types/dose";
import { configureNotifications, scheduleDoseNotification, cancelDoseNotification, cancelAllDoseNotifications, scheduleUpcomingDoseNotifications,  } from "../utils/notifications";
import { cacheGet, cacheSet } from "@/utils/cache";
import WeekCalendar from "@/components/WeekCalender";
import { format } from "date-fns";


export default function TodayScreen() {
  const [doses, setDoses] = useState<Dose[]>([]);
  const [selectedDate, setSelectedDate] = useState(new Date());
  const TODAY_CACHE_KEY = "today_doses";
  const UPCOMING_CACHE_KEY = "upcoming_doses";

  useEffect(() => {
    load();
  }, []);

  const load = async () => {
    const data = await getUpcomingDoses();
    setDoses(data);
  };

  const filtered = doses.filter(
    (d) =>
      format(new Date(d.scheduled_at), "yyyy-MM-dd") ===
      format(selectedDate, "yyyy-MM-dd")
  );

  useEffect(() => {
    const init = async () => {
      // 1️⃣ Load cached today doses
      const cached = await cacheGet<Dose[]>(TODAY_CACHE_KEY);
      if (cached) setDoses(cached);

      // 2️⃣ Fetch fresh today doses
      try {
        const upcomingData = await getUpcomingDoses();

        setDoses(upcomingData); 
        await cacheSet(UPCOMING_CACHE_KEY, upcomingData);

        // 3️⃣ Schedule notifications for today's pending doses
        // await Promise.all(
        //   todayData.filter(d => d.status === "pending").map(scheduleDoseNotification)
        // );

        await cancelAllDoseNotifications();
        await scheduleUpcomingDoseNotifications(upcomingData);
      } catch (e) {
        console.warn("Failed to fetch today doses, using cache", e);
      }

      // 4️⃣ Ensure upcoming notifications are scheduled
      try {
        const allowed = await configureNotifications();
        if (!allowed) return;

        const upcomingData = await getUpcomingDoses();
        await cacheSet(UPCOMING_CACHE_KEY, upcomingData);
        await scheduleUpcomingDoseNotifications(upcomingData);
      } catch (e) {
        console.warn("Failed to schedule upcoming notifications", e);
      }
    };

    init();
  }, []);

  useEffect(() => {
    const initNotifications = async () => {
      const allowed = await configureNotifications();
      if (!allowed) return;

      const upcoming = await getUpcomingDoses();

      await cancelAllDoseNotifications();
      await scheduleUpcomingDoseNotifications(upcoming);
    };

    initNotifications();
  }, []);


  const groupByTime = (doses: Dose[]) => {
    const map: Record<string, Dose[]> = {};

    doses.forEach((dose) => {
      const time = new Date(dose.scheduled_at).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });

      if (!map[time]) map[time] = [];
        map[time].push(dose);
    });

    return Object.entries(map).sort(([a], [b]) => {
      return a.localeCompare(b);
    });
  };

  const groups = groupByTime(filtered);

  return (
    <View style={{ padding: 16 }}>
      <WeekCalendar selected={selectedDate} onSelect={setSelectedDate} />

      <Text style={{ fontSize: 18, marginBottom: 12 }}>
        {format(selectedDate, "MMM d")}
      </Text>

      {groups.map(([time, doses]) => (
      <View key={time} style={{ marginBottom: 20 }}>
        <Text style={{ fontSize: 18, fontWeight: "600", marginBottom: 10 }}>
          {time}
        </Text>

        {doses.map((dose) => (
          <DoseCard
            key={dose.id}
            dose={dose}
            onTake={async () => {
                    await takeDose(dose.id);
                    await cancelDoseNotification(dose.id);
                    load();
                }}
            onSkip={async () => {
                    await skipDose(dose.id);
                    await cancelDoseNotification(dose.id);
                load();
                }}
          />
        ))}
      </View>
    ))}

      <FlatList
        data={filtered}
        keyExtractor={(d) => d.id}
        renderItem={({ item }) => (
          <DoseCard
            dose={item}
                onTake={async () => {
                    await takeDose(item.id);
                    await cancelDoseNotification(item.id);
                    load();
                }}

                onSkip={async () => {
                    await skipDose(item.id);
                    await cancelDoseNotification(item.id);
                load();
                }}
          />
        )}
      />
    </View>
  );
}

