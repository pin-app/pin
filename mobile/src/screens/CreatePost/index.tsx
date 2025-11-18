import React, { useCallback, useEffect, useRef, useState } from 'react';
import { View, Text, StyleSheet, SafeAreaView, TextInput, ScrollView, TouchableOpacity, Image, Alert, ActivityIndicator } from 'react-native';
import * as ImagePicker from 'expo-image-picker';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import Button from '@/components/Button';
import { apiService, Place } from '@/services/api';

interface SelectedImage {
  uri: string;
  remoteUrl?: string;
  uploading: boolean;
}

export default function CreatePostScreen() {
  const [placeQuery, setPlaceQuery] = useState('');
  const [placeResults, setPlaceResults] = useState<Place[]>([]);
  const [placeSuggestions, setPlaceSuggestions] = useState<Place[]>([]);
  const [placeMessage, setPlaceMessage] = useState<string | null>(null);
  const [selectedPlace, setSelectedPlace] = useState<Place | null>(null);
  const [caption, setCaption] = useState('');
  const [images, setImages] = useState<SelectedImage[]>([]);
  const [isSearchingPlaces, setIsSearchingPlaces] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const searchTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);
  const hasQuery = placeQuery.trim().length > 0;
  const displayedPlaces = hasQuery ? placeResults : placeSuggestions;
  const shouldShowPlaceList = !selectedPlace && displayedPlaces.length > 0;

  const loadSuggestedPlaces = useCallback(async () => {
    try {
      const suggestions = await apiService.listPlaces(8, 0);
      setPlaceSuggestions(suggestions);
    } catch (error) {
      console.error('Failed to load place suggestions', error);
    }
  }, []);

  useEffect(() => {
    loadSuggestedPlaces();
    return () => {
      if (searchTimeout.current) {
        clearTimeout(searchTimeout.current);
      }
    };
  }, [loadSuggestedPlaces]);

  const performSearch = async (query: string) => {
    try {
      setIsSearchingPlaces(true);
      const results = await apiService.searchPlaces(query, 8, 0);
      setPlaceResults(results);
      setPlaceMessage(results.length === 0 ? 'no places match that search' : null);
    } catch (error) {
      console.error('Failed to search places', error);
      setPlaceResults([]);
      setPlaceMessage('failed to search places');
    } finally {
      setIsSearchingPlaces(false);
    }
  };

  const handleSearchPlaces = (query: string) => {
    setPlaceQuery(query);
    setPlaceMessage(null);

    if (selectedPlace && query.trim() !== selectedPlace.name) {
      setSelectedPlace(null);
    }

    if (searchTimeout.current) {
      clearTimeout(searchTimeout.current);
    }

    if (!query.trim()) {
      setPlaceResults([]);
      return;
    }

    searchTimeout.current = setTimeout(() => performSearch(query.trim()), 250);
  };

  const handleSelectPlace = (place: Place) => {
    setSelectedPlace(place);
    setPlaceQuery(place.name);
    setPlaceResults([]);
    setPlaceMessage(null);
  };

  const handlePickImages = async () => {
    const { status } = await ImagePicker.requestMediaLibraryPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert('Permission needed', 'Allow photo library access to upload images.');
      return;
    }

    const result = await ImagePicker.launchImageLibraryAsync({
      allowsMultipleSelection: true,
      quality: 0.8,
      mediaTypes: ImagePicker.MediaTypeOptions.Images,
    });

    if (result.canceled) {
      return;
    }

    const newImages: SelectedImage[] = result.assets.map(asset => ({
      uri: asset.uri,
      uploading: true,
    }));

    setImages(prev => [...prev, ...newImages]);
    for (const asset of result.assets) {
      try {
        const uploadResult = await apiService.uploadPostImage(asset.uri, asset.fileName || undefined);
        setImages(prev =>
          prev.map(image =>
            image.uri === asset.uri && image.uploading
              ? { ...image, remoteUrl: uploadResult.url, uploading: false }
              : image
          )
        );
      } catch (error) {
        console.error('Failed to upload image', error);
        setImages(prev => prev.filter(image => image.uri !== asset.uri));
        Alert.alert('Upload failed', 'Could not upload one of the images.');
      }
    }
  };

  const handleRemoveImage = (uri: string) => {
    setImages(prev => prev.filter(image => image.uri !== uri));
  };

  const handleSubmit = async () => {
    if (!selectedPlace) {
      Alert.alert('Select a place', 'Choose where this post belongs.');
      return;
    }
    if (!caption.trim()) {
      Alert.alert('Add a caption', 'Share something about this post.');
      return;
    }
    const readyImages = images.filter(img => !!img.remoteUrl).map(img => img.remoteUrl!) ;
    if (!readyImages.length) {
      Alert.alert('Add images', 'Pick and upload at least one photo.');
      return;
    }

    try {
      setIsSubmitting(true);
      await apiService.createPost({
        place_id: selectedPlace.id,
        description: caption.trim(),
        images: readyImages,
      });
      setCaption('');
      setImages([]);
      setSelectedPlace(null);
      setPlaceQuery('');
      setPlaceResults([]);
      setPlaceMessage(null);
      loadSuggestedPlaces();
      Alert.alert('Posted', 'Your post is live.');
    } catch (error) {
      console.error('Failed to publish post', error);
      Alert.alert('Error', 'Could not publish this post.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.content} keyboardShouldPersistTaps="handled">
        <View style={styles.section}>
          <Text style={styles.label}>place</Text>
          <TextInput
            style={styles.input}
            placeholder="search places"
            placeholderTextColor={colors.textSecondary}
            value={placeQuery}
            onChangeText={handleSearchPlaces}
          />
          {isSearchingPlaces && (
            <ActivityIndicator color={colors.textSecondary} style={styles.loader} />
          )}
          {shouldShowPlaceList && (
            <View style={styles.results}>
              {displayedPlaces.map(place => (
                <TouchableOpacity
                  key={place.id}
                  style={styles.resultRow}
                  onPress={() => handleSelectPlace(place)}
                >
                  <Text style={styles.resultTitle}>{place.name}</Text>
                  <FontAwesome6 name="arrow-right" size={14} color={colors.textSecondary} />
                </TouchableOpacity>
              ))}
            </View>
          )}
          {!isSearchingPlaces && placeMessage && hasQuery && displayedPlaces.length === 0 && (
            <Text style={styles.hintText}>{placeMessage}</Text>
          )}
          {selectedPlace && (
            <View style={styles.selectedChip}>
              <Text style={styles.selectedText}>{selectedPlace.name}</Text>
              <TouchableOpacity onPress={() => setSelectedPlace(null)}>
                <FontAwesome6 name="xmark" size={14} color={colors.text} />
              </TouchableOpacity>
            </View>
          )}
        </View>

        <View style={styles.section}>
          <Text style={styles.label}>caption</Text>
          <TextInput
            style={[styles.input, styles.captionInput]}
            placeholder="how was it?"
            placeholderTextColor={colors.textSecondary}
            multiline
            value={caption}
            onChangeText={setCaption}
            maxLength={2000}
          />
          <Text style={styles.captionCount}>{caption.length}/2000</Text>
        </View>

        <View style={styles.section}>
          <View style={styles.imagesHeader}>
            <Text style={styles.label}>images</Text>
            <Button title="add" onPress={handlePickImages} variant="secondary" size="sm" />
          </View>
          <ScrollView horizontal showsHorizontalScrollIndicator={false}>
            {images.map(image => (
              <View key={image.uri} style={styles.imageWrapper}>
                <Image source={{ uri: image.uri }} style={styles.image} />
                {image.uploading && (
                  <View style={styles.imageOverlay}>
                    <ActivityIndicator color={colors.background} />
                  </View>
                )}
                <TouchableOpacity style={styles.removeButton} onPress={() => handleRemoveImage(image.uri)}>
                  <FontAwesome6 name="xmark" size={16} color={colors.background} />
                </TouchableOpacity>
              </View>
            ))}
          </ScrollView>
        </View>

        <Button
          title={isSubmitting ? 'posting...' : 'post'}
          onPress={handleSubmit}
          disabled={isSubmitting}
          style={styles.submit}
        />
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    padding: spacing.lg,
    gap: spacing.lg,
  },
  section: {
    borderWidth: 1,
    borderColor: colors.border,
    padding: spacing.md,
    borderRadius: 12,
    backgroundColor: colors.postBackground,
  },
  label: {
    fontSize: typography.fontSize.sm,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
    textTransform: 'uppercase',
  },
  input: {
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    color: colors.text,
    backgroundColor: colors.background,
  },
  loader: {
    marginTop: spacing.sm,
  },
  results: {
    marginTop: spacing.sm,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    maxHeight: 220,
    overflow: 'hidden',
  },
  resultRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.border,
  },
  resultTitle: {
    color: colors.text,
  },
  selectedChip: {
    marginTop: spacing.sm,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
  },
  selectedText: {
    color: colors.text,
  },
  hintText: {
    marginTop: spacing.xs,
    color: colors.textSecondary,
    fontSize: typography.fontSize.sm,
  },
  captionInput: {
    minHeight: 120,
    textAlignVertical: 'top',
  },
  captionCount: {
    textAlign: 'right',
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  imagesHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  imageWrapper: {
    width: 140,
    height: 180,
    borderRadius: 12,
    overflow: 'hidden',
    marginRight: spacing.md,
    backgroundColor: colors.border,
  },
  image: {
    width: '100%',
    height: '100%',
  },
  imageOverlay: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(0,0,0,0.5)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  removeButton: {
    position: 'absolute',
    top: spacing.xs,
    right: spacing.xs,
    backgroundColor: 'rgba(0,0,0,0.6)',
    borderRadius: 12,
    padding: 4,
  },
  submit: {
    marginTop: spacing.md,
  },
});

