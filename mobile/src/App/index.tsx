import React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { SafeAreaProvider } from "react-native-safe-area-context";
import TabNavigator from "./TabNavigator";
import { EventProvider } from "react-native-outside-press";

export { default as TabNavigator } from "./TabNavigator";

export default function App() {
  return (
    <SafeAreaProvider>
      <NavigationContainer>
        <EventProvider>
          <TabNavigator />
        </EventProvider>
      </NavigationContainer>
    </SafeAreaProvider>
  );
}
