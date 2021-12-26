import { useMutation, useQueryClient } from "react-query";
import { toast } from "react-toastify";
import { useAxios } from "../components/AuthContext";

export interface ITodoProps {
  Id: string;
  Completed: boolean;
  Content: string;
}

const Todo = (props: ITodoProps) => {
  const queryClient = useQueryClient();
  const axios = useAxios();
  const completeTodoMutation = useMutation(
    "completeTodo",
    (todoId: string) => axios.patch(`/todos/${todoId}`),
    {
      onSuccess: () => {
        queryClient.invalidateQueries(["todos"]);
        toast("Marked Complete", {
          style: { color: "green" },
          position: "top-right",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
          progressStyle: { background: '#0A5ECD' }
        });
      },
      onError: () => {
        toast("Could not mark Complete", {
          style: { color: "red" },
          position: "top-right",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
          progress: undefined,
          progressStyle: { background: '#0A5ECD' }
        });
      },
    }
  );

  const handleComplete = () => {
    completeTodoMutation.mutateAsync(props.Id);
  };

  return (
    <div
      key={props.Id}
      className="flex justify-between items-center px-4 py-2 border-b-2 last:border-b-0"
    >
      <div className="border-r-2 pr-4 flex-grow text-left">{props.Content}</div>
      <input
        type="checkbox"
        className="mr-4 ml-10 cursor-pointer"
        checked={props.Completed}
        onChange={handleComplete}
      />
    </div>
  );
};

export default Todo;
