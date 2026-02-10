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

  if (!["granted", "provisional"].includes(finalStatus)) {
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
  name: string;
  dosage: string;
  scheduled_at: string;
}) => {
  const triggerDate = new Date(dose.scheduled_at);

  if (triggerDate.getTime() <= Date.now()) return;

  await Notifications.scheduleNotificationAsync({
    identifier: dose.id,
    content: {
      title: "Medicine Reminder 💊",
      body: `Time to take ${dose.name} - ${dose.dosage} `,
      sound: "default",
    },
    trigger: {
      type: Notifications.SchedulableTriggerInputTypes.DATE,
      date: triggerDate,
    },
    channelId: Platform.OS === "android" ? "medicine" : undefined,  
  } as any);
};

export const cancelDoseNotification = async (doseId: string) => {
  await Notifications.cancelScheduledNotificationAsync(doseId);
};


export const cancelAllDoseNotifications = async () => {
  await Notifications.cancelAllScheduledNotificationsAsync();
};

export const scheduleUpcomingDoseNotifications = async (doses: Dose[]) => {
  await Promise.all(
  doses.map(async (dose) => {
    if (dose.status !== "pending") return;
    const triggerDate = new Date(dose.scheduled_at);
    if (triggerDate.getTime() <= Date.now()) return;
    await Notifications.cancelScheduledNotificationAsync(dose.id);
    await Notifications.scheduleNotificationAsync({
      identifier: dose.id,
      content: {
        title: "Medicine Reminder 💊",
        body: `Time to take ${dose.name} - ${dose.dosage}`,
        sound: "default",
      },
      trigger: {
        type: Notifications.SchedulableTriggerInputTypes.DATE,
        date: triggerDate,
      },
      channelId: Platform.OS === "android" ? "medicine" : undefined,
    } as any);
  })
);
};

