package ranking

import (
	"log"
)

// RankDocuments ranks the documents based on the query text
func RankDocuments(query Query) ([]Document, error) {
	// Get invertible index for the query
	index, err := getInvertibleIndex(query)
	if err != nil {
		return nil, err
	}

	// Get slice of all relevant documents
	documents, err := GetDocuments(index)
	if err != nil {
		return nil, err
	}

	// Count and avg length of all documents
	docStatistics := fetchTotalDocStatistics()

	// Add document metadata and features
	for _, document := range documents {
		document.Metadata, err = fetchDocumentMetadata(document.DocID)
		if err != nil {
			return nil, err
		}

		err = document.ComputeFeatures(query, docStatistics)
		if err != nil {
			return nil, err
		}
	}

	// add each feature to the feature struct
	// sort
	// rank

	log.Printf("Ranked documents for query: %s", query)

	// Return the ranked documents
	return documents, nil
}

// GetDocuments returns a slice of all documents in the InvertibleIndex
func GetDocuments(index InvertibleIndex) ([]Document, error) {
	// Map to store aggregated term frequencies for each document
	documentsMap := make(map[string]Document)

	// Iterate over each term in the invertible index
	for term, docIndices := range index {
		for _, docIndex := range docIndices {
			// Check if the document already exists in the map
			doc, exists := documentsMap[docIndex.DocID]
			if !exists {
				// If the document doesn't exist, initialize it
				doc = Document{
					DocID:           docIndex.DocID,
					TermFrequencies: make(map[string]int),
				}
			}

			// Add the term frequency to the document
			doc.TermFrequencies[term] += docIndex.Frequency
			documentsMap[docIndex.DocID] = doc
		}
	}

	// Convert the map to a slice
	documents := make([]Document, 0, len(documentsMap))
	for _, doc := range documentsMap {
		documents = append(documents, doc)
	}

	return documents, nil
}
