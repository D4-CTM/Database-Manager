import axios, { AxiosResponse, HttpStatusCode } from "axios"
const defaultErr = new Error('Network or unknown error')

export async function PostOrPut<T, t>(endpoint: string, data: t, isNew: boolean = false) {
    try {
        let response: AxiosResponse<T> = isNew
            ? await axios.post<T>(endpoint, data)
            : await axios.put<T>(endpoint, data)

        if (response.status === HttpStatusCode.Ok) {
            return response.data
        }
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}

export async function Get<T>(endpoint: string) {
    try {
        const response = await axios.get<T>(endpoint)

        if (response.status === HttpStatusCode.Ok) {
            return response.data
        }
    } catch (ex) {
        if (axios.isAxiosError(ex) && ex.response) {
            const err = ex.response.data as string 
            throw new Error(err)
        }

        console.error(defaultErr.message, ex)
        throw defaultErr
    }
}
