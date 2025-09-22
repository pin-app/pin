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
import { colors } from "@/theme/colors";
import { typography } from "@/theme/typography";
import OutsidePressHandler from "react-native-outside-press";

interface SearchBarProps {
  placeholder: string;
  value: string;
  onInputChange: (text: string) => void;
  onSearchPress: () => void;
  onClear: () => void;
  onFocus?: () => void;
  onBlur?: () => void;
  style?: StyleProp<ViewStyle>;
}

export default function SearchBar({
  placeholder,
  value,
  onInputChange,
  onSearchPress,
  onClear,
  onFocus,
  onBlur,
  style,
}: SearchBarProps) {
  const [isInputFocused, setIsInputFocused] = useState(false);

  const handleFocus = () => {
    setIsInputFocused(true);
    onFocus?.();
  };

  const handleBlur = () => {
    setIsInputFocused(false);
    onBlur?.();
  };

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
          onFocus={handleFocus}
          onBlur={handleBlur}
          onSubmitEditing={onSearchPress}
          returnKeyType="search"
          autoComplete="off"
          autoCorrect={false}
          autoCapitalize="none"
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
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
    paddingHorizontal: 16,
    paddingVertical: 8,
    backgroundColor: colors.searchBackground,
    borderRadius: 8,
    marginHorizontal: 16,
    marginBottom: 8,
  },
  containerFocused: {
    backgroundColor: colors.searchBackground,
  },
  textInput: {
    flex: 1,
    fontSize: typography.fontSize.base,
    color: colors.text,
  },
  clearButton: {
    height: "auto",
    width: "auto",
  },
});
