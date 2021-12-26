import Head from "next/head";
import Router, { useRouter } from "next/router";
import {
  Dispatch,
  FormEvent,
  SetStateAction,
  useEffect,
  useState,
} from "react";
import { useMutation } from "react-query";
import { authenticate, LoginRequestPayload } from "../api/auth";
import Loader from "../components/Loader";

export default function Home() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const loginMutation = useMutation(
    "authenticate",
    (formData: LoginRequestPayload) => authenticate(formData)
  );
  const isLoading = loginMutation.isLoading;
  const token = loginMutation.data?.token;

  const handleChange = (
    setter: Dispatch<SetStateAction<string>>,
    value: string
  ) => {
    setter(value);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    loginMutation.mutateAsync({ username, password });
  };

  useEffect(() => {
    if (token) {
      localStorage.setItem("token", token);
      router.push("/todos");
    }
  }, [token]);

  if (isLoading || token) {
    console.log("Loading");
    return <Loader />;
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen py-2">
      <Head>
        <title>Go - Todo App</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <h1 className="text-5xl text-bold">Login Page</h1>
      <main className="flex flex-col items-center justify-center w-full flex-1 px-20">
        <div className="w-full max-w-xl">
          <form
            className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
            onSubmit={handleSubmit}
          >
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="username"
              >
                Username
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="username"
                type="text"
                value={username}
                onChange={(e) => handleChange(setUsername, e.target.value)}
                placeholder="Username"
              />
            </div>
            <div className="mb-6">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="password"
              >
                Password
              </label>
              <input
                className="shadow appearance-none rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                id="password"
                type="password"
                value={password}
                onChange={(e) => handleChange(setPassword, e.target.value)}
                placeholder="******************"
              />
            </div>
            <div className="flex flex-col items-center justify-between">
              <button
                className="bg-blue-500 hover:bg-blue-700 text-white py-2 px-4 focus:outline-none focus:shadow-outline w-full rounded-lg font-semibold text-lg"
                type="submit"
              >
                Sign In
              </button>
            </div>
          </form>
          <button
            type="button"
            className="w-full bg-green-400 hover:bg-green-700 text-white my-2 p-3 rounded-lg font-semibold text-lg"
            onClick={() => router.push("/register")}
          >
            Create New Account
          </button>
          <p className="text-center text-gray-500 text-xs">
            &copy;2021 All rights reserved.
          </p>
        </div>
      </main>

      <footer className="flex items-center justify-center w-full h-24 border-t">
        <a
          className="flex items-center justify-center"
          href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          Made with <span className="text-red-500 px-2 text-2xl">♥</span> By
          Tejas
        </a>
      </footer>
    </div>
  );
}
