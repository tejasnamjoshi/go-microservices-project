import { useRouter } from "next/router";
import { FormEvent, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { useAxios } from "../components/AuthContext";
import Loader from "../components/Loader";
import Todo from "../components/Todo";
import { toast } from "react-toastify";
import Logout from "../components/Logout";

export interface ITodo {
  Id: string;
  Completed: boolean;
  Content: string;
}

const Todos = () => {
  const [todoContent, setTodoContent] = useState("");
  const token =
    typeof window !== "undefined" ? localStorage.getItem("token") : undefined;
  const axios = useAxios();
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data, isLoading, isError, refetch } = useQuery(
    "todos",
    () => axios.get<ITodo[]>("/todos"),
    {
      staleTime: 0,
    }
  );

  const createTodoMutation = useMutation(
    "addTodo",
    (formData: { content: string }) => axios.post("/todos", formData),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(["todos"]);
        toast("Todo Added", {
          style: { color: "green" },
          position: "top-right",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
          progressStyle: { background: "#0A5ECD" },
        });
      },
      onError: () => {
        toast("Todo could not be added", {
          style: { color: "red" },
          position: "top-right",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
          progressStyle: { background: "#0A5ECD" },
        });
      },
    }
  );

  if (isLoading) return <Loader />;

  if (!token || isError) {
    router.push("/");
    return null;
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setTodoContent("");
    createTodoMutation.mutateAsync({ content: todoContent });
  };

  return (
    <div className="text-center min-h-screen py-2 max-w-2xl m-auto">
      <div className="flex justify-end">
        <Logout />
      </div>
      <h1 className="text-3xl font-bold pb-10">List of Todos</h1>
      <form className="flex py-4" onSubmit={handleSubmit}>
        <input
          className="shadow appearance-none rounded text-gray-700 leading-tight focus:outline-none focus:shadow-outline flex-grow mr-2"
          type="text"
          value={todoContent}
          onChange={(e) => setTodoContent(e.target.value)}
        />
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none  focus:bg-blue-900 focus:shadow-inner ml-2"
        >
          Add Todo
        </button>
      </form>

      <div className="border-2 rounded-md">
        <>
          {data.data.map((todo) => {
            return <Todo key={todo.Id} {...todo} />;
          })}
        </>
      </div>
    </div>
  );
};

export default Todos;
