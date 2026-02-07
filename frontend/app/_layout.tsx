import AppNavigator from '@/navigation/AppNavigator';
import 'react-native-reanimated';


export const unstable_settings = {
  anchor: '(tabs)',
};

export default function RootLayout() {

  return (
      <AppNavigator />
  );
}
