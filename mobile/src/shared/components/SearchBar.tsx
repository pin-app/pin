import FontAwesome6 from "@expo/vector-icons/FontAwesome6";
import { useRef, useState } from "react";
import {
  Keyboard,
  Pressable,
  StyleProp,
  StyleSheet,
  TextInput,
  View,
  ViewStyle,
} from "react-native";
import { colors } from "../../theme/colors";
import { typography } from "../../theme/typography";
import OutsidePressHandler from "react-native-outside-press";

interface SearchBarProps {
  placeholder: string;
  value: string;
  onInputChange: (text: string) => void;
  onSearchPress: () => void;
  onClear: () => void;
  style?: StyleProp<ViewStyle>;
}

export default function SearchBar({
  placeholder,
  value,
  onInputChange,
  onSearchPress,
  onClear,
  style,
}: SearchBarProps) {
  const [isInputFocused, setIsInputFocused] = useState(false);

  return (
    <OutsidePressHandler onOutsidePress={() => Keyboard.dismiss()}>
      <View
        style={[
          styles.container,
          isInputFocused && styles.containerFocused,
          style,
        ]}
      >
        <FontAwesome6 name="magnifying-glass" size={12} color="black" />
        <TextInput
          placeholder={placeholder}
          value={value}
          onChangeText={onInputChange}
          style={styles.textInput}
          onFocus={() => setIsInputFocused(true)}
          onBlur={() => setIsInputFocused(false)}
          onSubmitEditing={onSearchPress}
          returnKeyType="search"
        />
        {!!value && (
          <Pressable onPress={onClear}>
            <FontAwesome6 name="xmark" size={12} color="black" />
          </Pressable>
        )}
      </View>
    </OutsidePressHandler>
  );
}

const styles = StyleSheet.create({
  container: {
    width: "100%",
    display: "flex",
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
    padding: 4,
    paddingHorizontal: 8,
    backgroundColor: "white",
    borderWidth: 1,
    borderRadius: 8,
    fontSize: typography.fontSize["2xl"],
    color: colors.text,
  },
  containerFocused: {
    outlineColor: "blue",
    outlineWidth: 1,
  },
  textInput: {
    flex: 1,
  },
  clearButton: {
    height: "auto",
    width: "auto",
  },
});
