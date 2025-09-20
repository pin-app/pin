import React from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors } from '@/theme/colors';
import { typography } from '@/theme/typography';

interface SearchResult {
  id: string;
  type: 'place' | 'member' | 'recent';
  title: string;
  subtitle?: string;
  icon?: string;
}

interface SearchResultsProps {
  results: SearchResult[];
  onResultPress: (result: SearchResult) => void;
  onClearRecent?: () => void;
  onClose?: () => void;
}

export default function SearchResults({ results, onResultPress, onClearRecent, onClose }: SearchResultsProps) {
  const renderResult = ({ item }: { item: SearchResult }) => (
    <TouchableOpacity
      style={styles.resultItem}
      onPress={() => onResultPress(item)}
    >
      <View style={styles.resultContent}>
        {item.icon && (
          <FontAwesome6 
            name={item.icon as any} 
            size={16} 
            color={colors.textSecondary} 
            style={styles.resultIcon}
          />
        )}
        <View style={styles.resultText}>
          <Text style={styles.resultTitle}>{item.title}</Text>
          {item.subtitle && (
            <Text style={styles.resultSubtitle}>{item.subtitle}</Text>
          )}
        </View>
      </View>
      {item.type === 'recent' && (
        <TouchableOpacity onPress={() => onClearRecent?.()}>
          <FontAwesome6 name="xmark" size={12} color={colors.textTertiary} />
        </TouchableOpacity>
      )}
    </TouchableOpacity>
  );

  return (
    <View style={styles.container}>
      <FlatList
        data={results}
        renderItem={renderResult}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        contentContainerStyle={styles.contentContainer}
        ListHeaderComponent={
          onClose ? (
            <View style={styles.header}>
              <TouchableOpacity onPress={onClose} style={styles.closeButton}>
                <FontAwesome6 name="xmark" size={16} color={colors.text} />
              </TouchableOpacity>
            </View>
          ) : null
        }
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 8,
  },
  closeButton: {
    padding: 4,
  },
  contentContainer: {
    paddingHorizontal: 16,
    paddingTop: 8,
  },
  resultItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  resultContent: {
    flexDirection: 'row',
    alignItems: 'center',
    flex: 1,
  },
  resultIcon: {
    marginRight: 12,
  },
  resultText: {
    flex: 1,
  },
  resultTitle: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    marginBottom: 2,
  },
  resultSubtitle: {
    fontSize: typography.fontSize.sm,
    color: colors.textSecondary,
  },
});
