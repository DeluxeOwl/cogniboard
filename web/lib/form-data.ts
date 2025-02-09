/**
 * Converts an object with a files array into a FormData compatible object with indexed file entries.
 * This is specifically designed to work with DTOs that have a `Files []huma.FormFile 'form:"files"'` field in Go.
 *
 * @param data The input data object that may contain a files array
 * @returns An object with the same fields as the input, but with files converted to indexed entries (files.0, files.1, etc.)
 */
export function convertFilesToFormDataFormat<T extends object>(
	data: T
): T & Record<`files.${number}`, File> {
	// Create a new object with the same properties as the input
	const formattedData = { ...data } as T & Record<`files.${number}`, File>;

	// Get the files array if it exists
	const files = (data as any).files;
	if (Array.isArray(files)) {
		// Remove the original files array
		delete (formattedData as any).files;

		// Add each file individually as files.0, files.1, etc.
		files.forEach((file: File, index: number) => {
			(formattedData as any)[`files.${index}`] = file;
		});
	}

	return formattedData;
}
