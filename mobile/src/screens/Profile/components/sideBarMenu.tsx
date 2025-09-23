import React, { useEffect, useMemo, useRef } from 'react';
import {
  Modal, View, Text, Pressable, StyleSheet, SafeAreaView,
  Animated, FlatList, Dimensions
} from 'react-native';

export type MenuItem = {
  key: string;
  label: string;
  subtitle?: string;
  destructive?: boolean;
  onPress?: () => void;
  hidden?: boolean;
};

type Props = {
  visible: boolean;
  onClose: () => void;
  title?: string;
  sections: { header?: string; items: MenuItem[] }[];
};

export default function SidebarMenu({ visible, onClose, title = 'Profile Menu', sections }: Props) {
  const { width: W } = Dimensions.get('window');
  const PANEL_W = Math.min(360, Math.round(W * 0.9));
  const START_X = PANEL_W + 16;

  const slideX = useRef(new Animated.Value(START_X)).current;
  const fade = useRef(new Animated.Value(0)).current;

  // freeze data during animation to avoid choppy re-renders
  const data = useMemo(() => sections, [sections]);

  useEffect(() => {
    if (visible) {
      fade.setValue(0);
      slideX.setValue(START_X);
      Animated.parallel([
        Animated.timing(fade, { toValue: 1, duration: 160, useNativeDriver: true }),
        Animated.spring(slideX, {
          toValue: 0,
          stiffness: 260,
          damping: 28,
          mass: 1,
          useNativeDriver: true,
        }),
      ]).start();
    } else {
      Animated.parallel([
        Animated.timing(fade, { toValue: 0, duration: 120, useNativeDriver: true }),
        Animated.spring(slideX, {
          toValue: START_X,
          stiffness: 240,
          damping: 26,
          mass: 1,
          useNativeDriver: true,
        }),
      ]).start();
    }
  }, [visible]);

  const renderRow = (mi: MenuItem, isLast: boolean) => {
    if (mi.hidden) return null;
    return (
      <Pressable
        key={mi.key}
        onPress={() => { onClose(); mi.onPress?.(); }}
        style={({ pressed }) => [styles.row, pressed && styles.rowPressed]}
        android_ripple={{ color: 'rgba(0,0,0,0.05)' }}
      >
        <View style={styles.rowText}>
          <Text numberOfLines={1} style={[styles.label, mi.destructive && styles.destructive]}>
            {mi.label}
          </Text>
          {!!mi.subtitle && !mi.destructive && (
            <Text numberOfLines={1} style={styles.subtitle}>{mi.subtitle}</Text>
          )}
        </View>
        {!mi.destructive && <Text style={styles.chev}>{'›'}</Text>}
        {!isLast && <View style={styles.separator} />}
      </Pressable>
    );
  };

  return (
    <Modal visible={visible} transparent animationType="none" onRequestClose={onClose}>
      <Animated.View style={[styles.backdrop, { opacity: fade }]}>
        <Pressable style={StyleSheet.absoluteFill} onPress={onClose} />
      </Animated.View>

      <Animated.View
        style={[
          styles.panel,
          { width: PANEL_W, transform: [{ translateX: slideX }], marginRight: 8 }
        ]}
      >
        <SafeAreaView style={styles.safe}>
          <View style={styles.headerRow}>
            <Text style={styles.title}>{title}</Text>
            <Pressable hitSlop={12} onPress={onClose}><Text style={styles.done}>Done</Text></Pressable>
          </View>

          <FlatList
            data={data}
            keyExtractor={(_, i) => `section-${i}`}
            contentContainerStyle={{ paddingBottom: 28 }}
            removeClippedSubviews
            windowSize={4}
            renderItem={({ item }) => {
              const items = item.items.filter(i => !i.hidden);
              if (!items.length) return null;
              return (
                <View style={styles.section}>
                  {item.header ? <Text style={styles.sectionHeader}>{item.header.toUpperCase()}</Text> : null}
                  <View style={styles.card}>
                    {items.map((mi, idx) => renderRow(mi, idx === items.length - 1))}
                  </View>
                </View>
              );
            }}
          />
        </SafeAreaView>
      </Animated.View>
    </Modal>
  );
}

const R = {
  radius: 22,
  padX: 20,
  spaceS: 12,
  spaceM: 16,
  hair: StyleSheet.hairlineWidth,
};

const styles = StyleSheet.create({
  backdrop: { ...StyleSheet.absoluteFillObject, 
    backgroundColor: 'rgba(0,0,0,0.25)'
  },

  panel: {
    position: 'absolute',
    right: -7,
    top: 0,
    bottom: 0,
    backgroundColor: '#fff',
    borderTopLeftRadius: 50,
    borderBottomLeftRadius: 20,
    shadowColor: '#000',
    shadowOpacity: 0.1,
    shadowRadius: 14,
    elevation: 6,
  },

  safe: {
    flex: 1,
    paddingHorizontal:
    R.padX
  },

  headerRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingTop: 6,
    paddingBottom: 8,
    left: 5
  },

  title: {
    fontSize: 17,
    fontWeight: '600',
    letterSpacing: -0.2
  },

  done: {
    fontSize: 16,
    fontWeight: '600',
    color: '#2563eb',
    right: 10
  },

  section: { 
    marginTop: 18

  },
  sectionHeader: {
    fontSize: 11,
    color: '#9ca3af',
    marginBottom: 8,
    letterSpacing: 0.6,
    left: 8
  },

  card: {
    borderRadius: 16,
    backgroundColor: '#fff',
    overflow: 'hidden',
    borderWidth: R.hair,
    borderColor: 'rgba(0,0,0,0.06)',
  },

  row: {
    paddingVertical: R.spaceS + 2,
    paddingHorizontal: R.padX,
    flexDirection: 'row',
    alignItems: 'center',
    minHeight: 46,
  },
  rowPressed: {
    backgroundColor: 'rgba(0,0,0,0.035)'
  },

  rowText: {
    flex: 1
  },

  label: {
    fontSize: 16,
    color: '#0f172a',
    fontWeight: '500'
  },

  subtitle: {
    fontSize: 12,
    color: '#6b7280',
    marginTop: 2
  },

  chev: {
    fontSize: 20,
    color: 'rgba(0,0,0,0.25)',
    marginLeft: 8
  },

  separator: {
    position: 'absolute',
    right: R.padX,
    left: R.padX,
    bottom: 0,
    height: R.hair,
    backgroundColor: 'rgba(0,0,0,0.06)'
  },
  destructive: {
    color: '#dc2626',
    fontWeight: '600'
  },
});
