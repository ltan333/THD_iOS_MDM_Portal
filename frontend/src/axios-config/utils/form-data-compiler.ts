type FormDataValue = string | number | boolean | File | Blob | null | undefined;

export interface FormDataParams {
  [key: string]: FormDataValue | FormDataValue[];
}

export const ObjectToFormData = (params: FormDataParams = {}): FormData => {
  const formData = new FormData();

  const appendValue = (key: string, value: FormDataValue) => {
    if (value === null || value === undefined) return;

    if (value instanceof File || value instanceof Blob) {
      formData.append(key, value);
    } else {
      formData.append(key, value.toString());
    }
  };

  for (const [key, value] of Object.entries(params)) {
    if (Array.isArray(value)) {
      value.forEach((item) => {
        appendValue(`${key}[]`, item);
      });
    } else {
      appendValue(key, value);
    }
  }

  return formData;
};
