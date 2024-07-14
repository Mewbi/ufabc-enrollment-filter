import axios from "axios";

type ErrorResponse = {
  error: boolean;
  message?: string;
};

type Enrollment = {
  id: string;
  name: string;
  create_at?: string;
};

type MayBeError<T> = T | ErrorResponse;

const API_URL = `${
  import.meta.env.PROD ? "ufabc-parser.felipef.software" : "localhost:8080"
}`;

const apiRequest = axios.create({
  baseURL: `${import.meta.env.PROD ? "https" : "http"}://${API_URL}`,
  withCredentials: true,
});

const getApiStatus = async (): Promise<boolean> => {
  const res = await apiRequest
    .get(`/health`)
    .then((response) => {
      if (response.status === 200) {
        return true;
      }

      return false;
    })
    .catch((_) => {
      return false;
    });

  return res;
};

const getEnrollmentList = async (): Promise<MayBeError<Enrollment[]>> => {
  const res = await apiRequest
    .get(`/enrollment`)
    .then((response) => response.data as Enrollment[])
    .catch((err) => {
      const data :ErrorResponse = {
        error: true
      };
      if (err.response && err.response.data.error)
        data.message = err.response.data.error;
      return data
  });

  return res;
};

const getEnrollment = async (id: string): Promise<MayBeError<string>> => {
  const res = await apiRequest
    .get(`/enrollment/${id}`)
    .then((response) => response.data as string)
    .catch((err) => {
      const data :ErrorResponse = {
        error: true
      };
      if (err.response && err.response.data.error)
        data.message = err.response.data.error;
      return data
  });

  return res;
};

const newEnrollment = async (name: string, url: string): Promise<MayBeError<Enrollment>> => {
  const res = await apiRequest
    .post(`/parse-enrollment`, {name, url})
    .then((response) => response.data as Enrollment)
    .catch((err) => {
      const data :ErrorResponse = {
        error: true
      };
      if (err.response && err.response.data.error)
        data.message = err.response.data.error;
      return data
  });

  return res;
};


const isErrorResponse = (value: any): value is ErrorResponse => {
  return (value as ErrorResponse).error !== undefined;
}

export default { getApiStatus, getEnrollmentList, getEnrollment, newEnrollment, isErrorResponse };
