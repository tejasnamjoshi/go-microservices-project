import { useRouter } from "next/router";
import { toast } from "react-toastify";

const Logout = () => {
  const router = useRouter();
  const handleLogout = () => {
    toast("You have been logged out", {
      style: { color: "red" },
      position: "top-right",
      autoClose: 3000,
      hideProgressBar: true,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
    });
    localStorage.clear();
    router.push("/");
  };

  return (
    <button
      type="button"
      className="text-white bg-red-700 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-full text-sm px-5 py-2.5 text-center mr-2 mb-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
      onClick={handleLogout}
    >
      Logout
    </button>
  );
};

export default Logout;
