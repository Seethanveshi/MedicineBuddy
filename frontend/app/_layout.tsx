import "../utils/backgroundTask";
import AppNavigator from '@/navigation/AppNavigator';
import 'react-native-reanimated';
import { cancelDoseNotification, configureNotifications, registerDoseActions } from '@/utils/notifications';
import { skipDose, takeDose } from '@/api/doses';
import { useEffect } from 'react';
import * as Notifications from "expo-notifications";

export const unstable_settings = {
  anchor: '(tabs)',
};

export default function RootLayout() {
  useEffect(() => {
    configureNotifications();
    registerDoseActions();
  }, []); 


  useEffect(() => {
    // Notifications.registerTaskAsync("doseActionTask");
    const subscription =
      Notifications.addNotificationResponseReceivedListener(
        async (response) => {
          const action = response.actionIdentifier;
          const doseId =
            response.notification.request.content.data?.doseId;

          console.log("response" , doseId);
          if (!doseId) return;

          try {
            if (action === "take") {
              await takeDose(doseId as string);
            }

            if (action === "skip") {
              await skipDose(doseId as string);
            }

            await Notifications.dismissNotificationAsync(
              response.notification.request.identifier
            );

            await cancelDoseNotification(doseId as string);
            console.log("responsed" , doseId);

          } catch (e) {
            console.warn("Action failed", e);
          }
        }
      );

    return () => subscription.remove();
  }, []);



  return (
    // <GestureHandlerRootView style={{ flex: 1 }}>
      <AppNavigator />
    /* </GestureHandlerRootView> */
  );
}
