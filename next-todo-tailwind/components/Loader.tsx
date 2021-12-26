const Loader = () => {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="w-16 h-16 border-4 border-blue-400 border-solid rounded-full animate-spin border-t-transparent">
          <span className="sr-only">Loading...</span>
      </div>
    </div>
  );
};

export default Loader;
