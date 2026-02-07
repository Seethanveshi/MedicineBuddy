import { Dose } from "@/types/dose";
import * as Notifications from "expo-notifications";
import { Platform } from "react-native";

export const configureNotifications = async () => {
  const { status } = await Notifications.getPermissionsAsync();

  let finalStatus = status;
  if (status !== "granted") {
    const req = await Notifications.requestPermissionsAsync();
    finalStatus = req.status;
  }

  if (finalStatus !== "granted") {
    console.warn("Notification permission not granted");
    return false;
  }

  if (Platform.OS === "android") {
    await Notifications.setNotificationChannelAsync("medicine", {
      name: "Medicine Reminders",
      importance: Notifications.AndroidImportance.MAX,
      vibrationPattern: [0, 250, 250, 250],
      sound: "default",
    });
  }

  return true;
};



export const scheduleDoseNotification = async (dose: {
  id: string;
  scheduled_at: string;
}) => {
  const triggerDate = new Date(dose.scheduled_at);

  // Ignore past dates
  if (triggerDate.getTime() <= Date.now()) return;

  await Notifications.scheduleNotificationAsync({
    identifier: dose.id,
    content: {
      title: "Medicine Reminder 💊",
      body: "Time to take your medicine",
      sound: "default",
    },
    trigger: {
      type: Notifications.SchedulableTriggerInputTypes.DATE,
      date: triggerDate,
    },
  });
};

export const cancelDoseNotification = async (doseId: string) => {
  await Notifications.cancelScheduledNotificationAsync(doseId);
};


export const cancelAllDoseNotifications = async () => {
  const scheduled = await Notifications.getAllScheduledNotificationsAsync();
  for (const n of scheduled) {
    if (n.identifier) {
      await Notifications.cancelScheduledNotificationAsync(n.identifier);
    }
  }
};

export const scheduleUpcomingDoseNotifications = async (doses: Dose[]) => {
  for (const dose of doses) {
    if (dose.status !== "pending") continue;

    const triggerDate = new Date(dose.scheduled_at);
    if (triggerDate.getTime() <= Date.now()) continue;

    await Notifications.scheduleNotificationAsync({
      identifier: dose.id,
      content: {
        title: "Medicine Reminder 💊",
        body: "Time to take your medicine",
        sound: "default",
      },
      trigger: {
        type: Notifications.SchedulableTriggerInputTypes.DATE,
        date: triggerDate,
    },
    });
  }
};

