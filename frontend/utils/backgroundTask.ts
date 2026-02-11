import * as TaskManager from 'expo-task-manager';
import * as Notifications from 'expo-notifications';
import { takeDose, skipDose } from '@/api/doses';
import { cancelDoseNotification } from '@/utils/notifications';

TaskManager.defineTask('doseActionTask', async () => {
  console.log("✅ Background task fired!");
});

TaskManager.defineTask('doseActionTask', async ({ data, error }) => {
  if (error) return;

  const { notification, actionIdentifier } = data as {
    notification: Notifications.Notification;
    actionIdentifier: string;
  };

  const doseId = notification.request.content.data?.doseId;
  if (!doseId) return; 

  try {
    if (actionIdentifier === 'take') await takeDose(doseId as string);
    if (actionIdentifier === 'skip') await skipDose(doseId as string);

    await cancelDoseNotification(doseId as string);
    console.log('Handled dose in background:', doseId);
  } catch (err) {
    console.error(err);
  }
});

// Register headless task
Notifications.registerTaskAsync('doseActionTask');
