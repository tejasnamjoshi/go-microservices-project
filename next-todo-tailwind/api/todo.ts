import axios from 'axios'

const baseUrl = 'http://localhost:3001'

export interface LoginRequestPayload {
    username: string;
    password: string;
}

export const authenticate = async (formData: LoginRequestPayload) => {
    const { data } = await axios.post(`${baseUrl}/login`, formData)
    return data
}