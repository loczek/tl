import * as z from "zod/mini";

export async function postData<T extends z.ZodMiniType>(
  url: string,
  data: object,
  responseSchema: T,
) {
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    throw new Error(`HTTP error status: ${response.status}`);
  }

  return responseSchema.parse(await response.json());
}

export async function getData<T extends z.ZodMiniType>(
  url: string,
  responseSchema: T,
) {
  const response = await fetch(url, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`HTTP error status: ${response.status}`);
  }

  return responseSchema.parse(await response.json());
}
