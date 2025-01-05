import { LogOut, Settings } from 'lucide-react';
import { useRouter } from 'next/navigation';
import React from 'react';

const options = [
  { icon: Settings, name: "Settings" },
  { icon: LogOut, name: "Sign Out" }
]

const Profile = ({ shown, close }: { shown: boolean, close: () => void, className?: string }) => {
  const router = useRouter();

  const handleSignOut = () => {
    localStorage.removeItem("token")
    router.push("/login")
  }

  return shown && (
    <div
      className='fixed z-[2] top-0 bottom-0 left-0 right-0 w-full h-full translate-all duration-200 bg-transparent'
      onClick={() => close()}
    >
      <div
        className="absolute top-20 right-6 bg-white dark:bg-neutral-900 p-3 rounded-lg flex flex-col gap-2"
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <div className='flex flex-col gap-3'>
          {options?.map((option, index) => (
            <div
              key={index}
              onClick={
                option.name === "Sign Out" ? (
                  handleSignOut
                ) : (
                  () => router.push("/profile")
                )
              }
              className='flex items-center gap-2 cursor-pointer'
            >
              <option.icon className={`${option.name === "Sign Out" && "text-red-500"}`} />
              <p>{option.name}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Profile;
